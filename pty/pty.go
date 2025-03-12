package pty

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/System233/ll-killer-go/utils"

	"github.com/containerd/console"
	"github.com/moby/sys/reexec"
	"github.com/moby/term"
)

type Pty struct {
	Socket          string
	AutoExit        bool
	Timeout         time.Duration
	connectionCount int32
	connectionTimes int32
}
type PtyExecArgs struct {
	Args   []string
	Pty    string
	Path   string
	Dir    string
	NsType []string
	// NoFail bool
	Env []string
}
type PtyExecReply struct {
	ExitCode int
}
type PtyNsInfoArgs struct {
}
type PtyNsInfoRelpy struct {
	Pid int
	Uid int
	Gid int
}

func (pty *Pty) NsInfo(args *PtyNsInfoArgs, reply *PtyNsInfoRelpy) error {
	reply.Pid = os.Getpid()
	reply.Gid = os.Getgid()
	reply.Uid = os.Getuid()
	return nil
}
func (pty *Pty) Exec(args *PtyExecArgs, reply *PtyExecReply) error {
	utils.Debug("RPC:Exec:", args)
	slave, err := os.OpenFile(args.Pty, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer slave.Close()

	cmd := exec.Command(args.Args[0], args.Args[1:]...)
	if args.Env != nil {
		cmd.Env = args.Env
	}
	if args.Path != "" {
		cmd.Path = args.Path
	}
	if args.Dir != "" {
		cmd.Dir = args.Dir
	}

	cmd.Stdin = slave
	cmd.Stdout = slave
	cmd.Stderr = slave
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setctty:   true,
		Setsid:    true,
		Pdeathsig: syscall.SIGKILL,
	}

	err = cmd.Start()
	if err != nil {
		utils.Debug("RPC:Exec:Err", err)
		return err
	} else {
		utils.Debug("RPC:Exec:Pid", cmd.Process.Pid)
	}
	pid := cmd.Process.Pid
	err = cmd.Wait()

	utils.Debug("RPC:Exec:Exited", pid, err)
	var exitErr *exec.ExitError
	if err != nil && errors.As(err, &exitErr) {
		reply.ExitCode = exitErr.ExitCode()
		return nil
	}
	return err
}
func (pty *Pty) Serve() error {
	if err := rpc.Register(pty); err != nil {
		utils.Debug("RPC:Register:", err)
		return err
	}

	if err := os.Remove(pty.Socket); err != nil && !os.IsNotExist(err) {
		utils.Debug("Server:remove", err)
		return err
	}
	listener, err := net.Listen("unix", pty.Socket)
	if err != nil {
		utils.Debug("RPC:Listen:", err)
		return err
	}
	defer listener.Close()
	utils.Debug("RPC:Listen on:", pty.Socket)

	var signal chan int = make(chan int)
	if pty.Timeout != 0 {
		go func() {
			time.Sleep(pty.Timeout)
			signal <- 1
		}()
	}
	go func() {
		for {
			conn, err := listener.Accept()
			atomic.AddInt32(&pty.connectionTimes, 1)
			if err != nil {
				signal <- 1
				utils.Debug("RPC:Accept:", err)
				continue
			}
			go func() {
				atomic.AddInt32(&pty.connectionCount, 1)
				rpc.ServeConn(conn)
				atomic.AddInt32(&pty.connectionCount, -1)
				signal <- 1
			}()
		}
	}()
	for range signal {
		if pty.AutoExit {
			if atomic.LoadInt32(&pty.connectionTimes) > 0 && atomic.LoadInt32(&pty.connectionCount) == 0 {
				utils.Debug("Server:Exit:AutoExit")
				return nil
			}
		}
		if pty.Timeout != 0 {
			if atomic.LoadInt32(&pty.connectionTimes) == 0 {
				utils.Debug("Server:Exit:Timeout")
				return nil
			}
		}
	}
	return nil
}

func (pty *Pty) connect() (*rpc.Client, error) {
	startTime := time.Now()
	for {
		conn, err := rpc.Dial("unix", pty.Socket)
		if err != nil {
			if pty.Timeout != 0 && time.Since(startTime) < pty.Timeout {
				time.Sleep(50 * time.Millisecond)
			} else {
				return nil, err
			}
		} else {
			return conn, nil
		}
	}
}
func (pty *Pty) Call(args *PtyExecArgs) (int, error) {

	client, err := pty.connect()
	if err != nil {
		utils.Debug("rpc.Dial:", err)
		return 1, err
	}
	defer client.Close()

	if args.Pty == "" {
		con, slavePath, err := console.NewPty()
		if err != nil {
			utils.Debug("console.NewPty:", err)
			return 1, err
		}
		defer con.Close()
		args.Pty = slavePath

		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			syscall.SIGWINCH,
			syscall.SIGINT,  // Ctrl+C 中断信号
			syscall.SIGTERM, // 终止信号
			syscall.SIGQUIT, // 退出信号
			syscall.SIGHUP,  // 挂断信号 / 重新加载配置
			syscall.SIGUSR1, // 用户自定义信号 1
			syscall.SIGUSR2, // 用户自定义信号 2
			syscall.SIGCONT, // 继续执行信号（当进程被暂停后恢复）
		)
		go func() {
			for sig := range ch {
				if sig == syscall.SIGWINCH {
					winSize, err := term.GetWinsize(os.Stdin.Fd())
					if err != nil {
						utils.Debug("term.GetSize:", err)
						continue
					}
					err = term.SetWinsize(con.Fd(), winSize)
					if err != nil {
						utils.Debug("term.SetWinsize:", err)
					}
				} else {
					var pgrp int
					_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, con.Fd(), uintptr(syscall.TIOCGPGRP), uintptr(unsafe.Pointer(&pgrp)))
					if errno != 0 {
						utils.Debug("syscall.Ioctl:", err)
						continue
					}
					err = syscall.Kill(-pgrp, sig.(syscall.Signal))
					if err != nil {
						utils.Debug("syscall.Kill:", err)
						continue
					}
				}
			}
		}()
		ch <- syscall.SIGWINCH

		if term.IsTerminal(os.Stdin.Fd()) {
			state, err := term.SetRawTerminal(os.Stdin.Fd())
			if err != nil {
				utils.Debug("term.SetRawTerminal:", err)
				return 1, err
			} else {
				defer term.RestoreTerminal(os.Stdin.Fd(), state)
			}
		}
		go func() {
			io.Copy(con, os.Stdin)
		}()
		go func() {
			io.Copy(os.Stdout, con)
		}()

	}

	utils.Debug("Pty.SelfPID", os.Getpid())
	var reply PtyExecReply
	if err := client.Call("Pty.Exec", args, &reply); err != nil {
		utils.Debug("rpc.Call:", err)
		return 1, err
	}
	return reply.ExitCode, nil
}

func (pty *Pty) NsEnter(execArgs *PtyExecArgs) (int, error) {
	utils.Debug("NsEnter", execArgs)
	client, err := pty.connect()
	if err != nil {
		utils.Debug("rpc.Dial:", err)
		return 1, err
	}
	defer client.Close()
	var args PtyNsInfoArgs
	var reply PtyNsInfoRelpy
	if err := client.Call("Pty.NsInfo", args, &reply); err != nil {
		utils.Debug("rpc.Pid:", err)
		return 1, err
	}
	utils.Debug("NsEnter.NsInfo", reply)
	cmdArgs := append([]string{"nsenter",
		fmt.Sprint(reply.Pid),
		"-kf",
		"-N", strings.Join(execArgs.NsType, ","),
		"--"}, execArgs.Args...)
	cmd := exec.Command(reexec.Self(), cmdArgs...)
	if execArgs.Env != nil {
		cmd.Env = execArgs.Env
	}
	if execArgs.Dir != "" {
		cmd.Dir = execArgs.Dir
	}
	if execArgs.Path != "" {
		cmd.Path = execArgs.Path
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			return exitErr.ExitCode(), nil
		}
		return 1, err
	}
	return 0, nil
}
