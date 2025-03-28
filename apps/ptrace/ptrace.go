/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _ptrace

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	internal "github.com/System233/ll-killer-go/apps/ptrace/internal"
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	seccomp "github.com/seccomp/libseccomp-golang"
)

const PtraceCommandHelp = `
此命令用于拦截并忽略系统调用错误，默认处理chown/fchown/fchownat/lchown系统调用。

seccomp模式下，可以忽略任意系统调用的错误（默认启用）。
ptrace模式下，仅处理上述系统调用，且仅支持amd64/arm64/loong64架构，其他架构下无效。

# 用法
<program> ptrace -- 要处理的命令
`

const IsSupported = internal.IsSupported

var Flag struct {
	Seccomp bool
	Syscall []string
}

func Ptrace(self string, args []string) error {
	args = append([]string{self, "ptrace", "--"}, args...)
	return utils.ExecRaw(args...)
}
func HandlePtraceEvent(process *os.Process, pid int) error {
	utils.Debug("HandlePtraceEvent", pid)
	var usage unix.Rusage
	var wstatus unix.WaitStatus
	var wpid int
	err := unix.Prctl(unix.PR_SET_PDEATHSIG, uintptr(unix.SIGKILL), 0, 0, 0)
	if err != nil {
		return err
	}
	err = unix.Prctl(unix.PR_SET_CHILD_SUBREAPER, 1, 0, 0, 0)
	if err != nil {
		return err
	}
	wpid, err = unix.Wait4(pid, &wstatus, unix.WUNTRACED, &usage)
	if err != nil {
		return err
	}
	err = unix.PtraceSetOptions(pid,
		unix.PTRACE_O_TRACECLONE|
			unix.PTRACE_O_TRACESYSGOOD|
			unix.PTRACE_O_TRACEFORK|
			unix.PTRACE_O_TRACEEXEC|
			unix.PTRACE_O_TRACEVFORK)
	if err != nil {
		return err
	}

	err = unix.PtraceSyscall(pid, 0)
	if err != nil {
		return err
	}
	IsError := func(err error) bool {
		if err == nil {
			return false
		}
		if err == unix.ESRCH || err == unix.ECHILD {
			return false
		}
		return true
	}
	for {
		wpid, err = unix.Wait4(-1, &wstatus, 0, &usage)
		if err != nil {
			return fmt.Errorf("Wait4:%v", err)
		}
		if wstatus.Exited() {
			/*
				If a thread group leader is traced and exits by calling _exit(2),
				a PTRACE_EVENT_EXIT stop will happen for it (if requested),
				but the subsequent WIFEXITED notification will not be delivered
				until all other threads exit. As explained above, if one of other
				threads calls execve(2), the death of the thread group leader
				will never be reported. If the execed thread is not traced by
				this tracer, the tracer will never know that execve(2) happened.
				One possible workaround is to PTRACE_DETACH the thread group
				leader instead of restarting it in this case.
				Last confirmed on 2.6.38.6.
			*/
			if wpid == process.Pid || process.Signal(syscall.Signal(0)) != nil {
				return &utils.ExitStatus{ExitCode: wstatus.ExitStatus()}
			}
			continue
		}
		if wstatus.Signaled() {
			if wpid == process.Pid || process.Signal(syscall.Signal(0)) != nil {
				return &utils.ExitStatus{ExitCode: -int(wstatus.Signal())}
			}
			continue
		}
		if wstatus.Continued() || !wstatus.Stopped() {
			continue
		}
		if wstatus.StopSignal()&0x80 == 0 {
			sig := wstatus.StopSignal()
			if wstatus.StopSignal() == unix.SIGTRAP || (wpid != pid && wstatus.StopSignal() == unix.SIGSTOP) {
				sig = 0
			}
			err = unix.PtraceSyscall(wpid, int(sig))
			if err != nil {
				return fmt.Errorf("PtraceSyscall.SIGTRAP:%#x,%v", wstatus, err)
			}
			continue
		}
		var regs unix.PtraceRegs
		err = unix.PtraceGetRegs(wpid, &regs)
		if IsError(err) {
			return fmt.Errorf("PtraceGetRegs: %#x, %v", wstatus, err)
		}
		err = internal.PtraceHandle(wpid, regs)
		if IsError(err) {
			return fmt.Errorf("PtraceHandle:%v", err)
		}
		err = unix.PtraceSyscall(wpid, 0)
		if IsError(err) {
			return fmt.Errorf("PtraceSyscall:%v", err)
		}
	}
}
func SetupSeccomp(ctx *seccomp.ScmpFilter, filters []string) error {
	ActIgnore := seccomp.ActErrno.SetReturnCode(0)

	for _, name := range filters {
		opcode, err := seccomp.GetSyscallFromName(name)
		if err != nil {
			utils.Debug("Seccomp: syscall not found", name)
			continue
		}
		utils.Debug("Seccomp: AddRule ActIgnore", name, opcode)
		if err := ctx.AddRule(opcode, ActIgnore); err != nil {
			return err
		}
	}
	return nil
}
func PtraceMain(cmd *cobra.Command, args []string) error {
	utils.Debug("PtraceMain", args)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if len(args) == 0 {
		args = []string{utils.DefaultShell()}
	}

	binary, err := exec.LookPath(args[0])
	if err == nil {
		args[0] = binary
	}

	if Flag.Seccomp {
		ctx, err := seccomp.NewFilter(seccomp.ActAllow)
		if err != nil {
			return err
		}
		defer ctx.Release()

		if err := SetupSeccomp(ctx, Flag.Syscall); err != nil {
			return err
		}

		if err := ctx.Load(); err != nil {
			return err
		}

		process, err := os.StartProcess(args[0], args, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Env:   os.Environ(),
		})
		if err != nil {
			return err
		}
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan,
			syscall.SIGWINCH,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGHUP,
			syscall.SIGUSR1,
			syscall.SIGUSR2,
			syscall.SIGCONT,
		)
		defer func() {
			signal.Stop(sigChan)
			close(sigChan)
		}()
		go func() {
			for sig := range sigChan {
				err := process.Signal(sig)
				if err != nil {
					utils.Debug("SignalError", err)
				}
			}
		}()
		status, err := process.Wait()
		if err != nil {
			return err
		}
		if status.ExitCode() != 0 {
			return &utils.ExitStatus{
				ExitCode: status.ExitCode(),
			}
		}
		return nil
	}

	process, err := os.StartProcess(args[0], args, &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Env:   os.Environ(),
		Sys: &syscall.SysProcAttr{
			Ptrace: true,
		},
	})

	if err != nil {
		return err
	}

	err = HandlePtraceEvent(process, process.Pid)
	if err != nil {
		return err
	}
	return nil
}
func CreatePtraceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ptrace",
		Short: "修正系统调用(chown)",
		Long:  utils.BuildHelpMessage(PtraceCommandHelp),
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(PtraceMain(cmd, args))
		},
	}
	cmd.Flags().BoolVar(&Flag.Seccomp, "seccomp", true, "使用seccomp而不是ptrace进行系统调用拦截")
	cmd.Flags().StringSliceVar(&Flag.Syscall, "syscall", []string{"chown", "fchown", "fchownat", "lchown"}, "指定拦截/忽略的系统调用")

	return cmd
}
