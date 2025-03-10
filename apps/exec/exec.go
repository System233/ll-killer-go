/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _exec

import (
	"fmt"
	"ll-killer/pty"
	"ll-killer/utils"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/moby/sys/reexec"
	"github.com/spf13/cobra"
)

var ExecFlag struct {
	Mounts []string
	// UidMappings   []string
	// GidMappings   []string
	UID               int
	GID               int
	Args              []string
	RootFS            string
	Socket            string
	SocketTimeout     time.Duration
	AutoExit          bool
	Root              bool
	NoFail            bool
	NoBindRootFS      bool
	FuseOverlayFS     string
	FuseOverlayFSArgs string
}

const ExecCommandHelp = `
 此命令在构建完成后的容器中作为入口点使用，无需在开发环境使用此命令。
 
 # 进入运行时环境，不挂载文件系统，也不进行 chroot，使用默认的用户和组
 <program> exec -- bash
 
 # 进入运行时环境，使用当前用户的 UID 和 GID，merge合并挂载文件系统，使用/myrootfs的内容覆盖根文件系统
 <program> exec --mount /+/myrootfs:/::merge -- bash
 
 # 进入运行时环境，挂载源路径到目标路径，使用指定的用户和组
 <program> exec --mount /path/to/source:/path/to/target --uid 2000 --gid 2000 --chroot=false
 
 # 使用自定义的 Unix 套接字和合并根目录路径
 <program> exec --socket /tmp/myapp.sock --rootfs /tmp/myapp.rootfs --mount /etc:/etc
 
 # 进入运行时环境，指定不同的挂载选项
 <program> exec --mount /data:/data:rw+nosuid --uid 1000 --gid 1000 --rootfs /var/run/myapp.rootfs
 
 ## 挂载选项
 使用参数 --mount 可以挂载文件或目录，参数与系统mount命令类似，但ll-killer额外支持merge合并挂载类型。
 用法：
	 --mount 源目录:目标目录:挂载标志:文件系统类型:额外数据
 
 挂载标志默认为bind或rbind，文件系统类型默认为none
 
 支持以下挂载标志：
 active async bind rbind dirsync invalidate i_version kernmount mandlock move noatime
 nodev nodiratime noexec nosuid nouser posixacl private rdonly rec relatime remount shared
 silent slave strictatime sync synchronous unbindable
 
 ### 合并挂载
 ll-killer额外支持merge合并挂载类型，用于在没有内核overlapfs或fuse模块支持的情况下堆叠文件系统。
 merge挂载将源目录列表中的文件从右到左堆叠，在目录存在冲突文件的情况下，越往后的目录优先级越高。
 若目录在所有源目录中只出现一次，那么该目录将直接绑定到目标目录中的相应位置，如果该目录只读，则挂载后的对应文件夹保持只读属性。
 用法：
	 --mount 源目录1+源目录2+源目录N:目标目录:挂载标志:merge:排除目录1+排除目录2+排除目录N
 
 源目录：支持多个源目录，多个源目录使用'+'分割。
 目标目录：合并文件系统的挂载位置。
 挂载标志：默认为bind/rbind
 文件系统：必须是merge
 排除目录：用于防止递归合并自身，或合并特殊文件系统。默认为: 目标目录+/proc+/dev+/tmp+/run+/var/run+/sys
		 排除目录将直接绑定到源目录中首次出现的路径
 
 
 * 理论上使用fuse模块实现效果最佳，但考虑到跨平台和发行版的问题，暂时使用合并挂载实现。
 
 `

func GetExecArgs() []string {
	args := []string{"--uid", fmt.Sprint(ExecFlag.UID), "--gid", fmt.Sprint(ExecFlag.GID)}

	if ExecFlag.FuseOverlayFS != "" {
		args = append(args, "--fuse-overlayfs", ExecFlag.FuseOverlayFS)
	}

	if ExecFlag.FuseOverlayFSArgs != "" {
		args = append(args, "--fuse-overlayfs-args", ExecFlag.FuseOverlayFSArgs)
	}

	if ExecFlag.RootFS != "" {
		args = append(args, "--rootfs", ExecFlag.RootFS)
	}

	if ExecFlag.Root {
		args = append(args, "--root")
	}

	if ExecFlag.NoFail {
		args = append(args, "--no-fail")
	}

	if ExecFlag.NoBindRootFS {
		args = append(args, "--no-bind-rootfs")
	}

	if ExecFlag.Socket != "" {
		args = append(args, "--socket", ExecFlag.Socket)
	}

	if ExecFlag.SocketTimeout != 0 {
		args = append(args, "--socket-timeout", fmt.Sprint(ExecFlag.SocketTimeout))
	}

	if !ExecFlag.AutoExit {
		args = append(args, "--auto-exit="+fmt.Sprint(ExecFlag.AutoExit))
	}

	for _, mount := range ExecFlag.Mounts {
		args = append(args, "--mount", mount)
	}

	if len(ExecFlag.Args) > 0 {
		args = append(args, "--")
		args = append(args, ExecFlag.Args...)
	}
	return args
}

func NewPtyFromFlags() *pty.Pty {
	return &pty.Pty{
		Socket:   ExecFlag.Socket,
		Timeout:  ExecFlag.SocketTimeout,
		AutoExit: ExecFlag.AutoExit,
	}
}
func PivotRootSystem() {
	utils.Debug("PivotRootSystem")
	if !ExecFlag.NoBindRootFS {
		err := utils.MountBind(ExecFlag.RootFS, ExecFlag.RootFS, 0)
		if err != nil {
			utils.ExitWith(err)
		}
	}
	oldRootFS := fmt.Sprint(ExecFlag.RootFS, ExecFlag.RootFS)
	utils.Debug("PivotRoot", ExecFlag.RootFS, oldRootFS)
	if err := syscall.PivotRoot(ExecFlag.RootFS, oldRootFS); err != nil {
		utils.ExitWith(err)
	}
	ExecShell()
}

func ExecSystem() {
	utils.Debug("ExecSystem")
	if ExecFlag.Socket != "" {
		pty := NewPtyFromFlags()
		pty.Serve()
	} else {
		utils.Exec(ExecFlag.Args...)
	}
}
func ExecShell() {
	if ExecFlag.UID == 0 && ExecFlag.GID == 0 || ExecFlag.Root {
		utils.Exec(ExecFlag.Args...)
		return
	}
	err := utils.SwitchTo("ExecSystem", &utils.SwitchFlags{
		UID:        ExecFlag.UID,
		GID:        ExecFlag.GID,
		Cloneflags: syscall.CLONE_NEWUSER,
	})
	if err != nil {
		utils.ExitWith(err)
	}
}
func MountFileSystem() {
	utils.Debug("MountFileSystem")
	isFuseOverlayFs := false
	for _, mount := range ExecFlag.Mounts {
		opt := utils.ParseMountOption(mount)
		err := opt.Mount()
		if err != nil {
			if ExecFlag.NoFail {
				utils.ExitWith(err, "mount", mount)
			}
			log.Println(err)
		}
		if opt.FSType == utils.FuseOverlayFSType {
			isFuseOverlayFs = true
		}
	}

	if ExecFlag.RootFS != "" {
		oldRootFS := fmt.Sprint(ExecFlag.RootFS, ExecFlag.RootFS)
		err := utils.MkdirAlls([]string{oldRootFS}, 0755)
		if err != nil {
			utils.ExitWith(err)
		}

		if isFuseOverlayFs {
			err = utils.SwitchTo("PivotRootSystem", &utils.SwitchFlags{Cloneflags: syscall.CLONE_NEWNS})
			if err != nil {
				utils.ExitWith(err)
			}
		} else {
			PivotRootSystem()
		}
	} else {
		ExecShell()
	}
}
func StartMountFileSystem() error {
	return utils.SwitchTo("MountFileSystem", &utils.SwitchFlags{
		UID:           0,
		GID:           0,
		Cloneflags:    syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Args:          append([]string{"exec"}, GetExecArgs()...),
		NoDefaultArgs: true,
	})
}
func ExecMain(cmd *cobra.Command, args []string) error {
	utils.Debug("ExecMain", args)
	ExecFlag.Args = args
	utils.GlobalFlag.FuseOverlayFS = ExecFlag.FuseOverlayFS
	utils.GlobalFlag.FuseOverlayFSArgs = ExecFlag.FuseOverlayFSArgs
	reexec.Register("MountFileSystem", MountFileSystem)
	reexec.Register("ExecSystem", ExecSystem)
	reexec.Register("PivotRootSystem", PivotRootSystem)
	if !reexec.Init() {
		if ExecFlag.Socket != "" {
			var signal chan error = make(chan error)
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			ptyFlag := NewPtyFromFlags()
			ptyFlag.Timeout = 0
			args := &pty.PtyExecArgs{
				Args: args,
				Dir:  cwd,
				Env:  os.Environ(),
			}
			exitCode, err := ptyFlag.Call(args)
			utils.Debug("pty.Call", exitCode, err)
			if err != nil {
				go func() {
					signal <- StartMountFileSystem()
				}()
				go func() {
					ptyFlag.Timeout = ExecFlag.SocketTimeout
					exitCode, err := ptyFlag.Call(args)
					if err == nil {
						err = &utils.ExitStatus{ExitCode: exitCode}
					}
					signal <- err
				}()
				return <-signal
			}
			return &utils.ExitStatus{ExitCode: exitCode}
		} else {
			return StartMountFileSystem()
		}
	}
	return nil
}

func CreateExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "进入运行时环境",
		Long:  utils.BuildHelpMessage(ExecCommandHelp),
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(ExecMain(cmd, args))
		},
	}

	// cmd.Flags().StringSliceVar(&ExecFlag.Mounts, "mount", []string{}, "source:target:[flags:[fstype:[option]]]")
	cmd.Flags().StringArrayVar(&ExecFlag.Mounts, "mount", []string{}, "source:target:[flags:[fstype:[option]]]")
	// cmd.Flags().StringArrayVar(&ExecFlag.UidMappings, "uidmapping", []string{}, "source:target:[flags:[fstype:[option]]]")
	// cmd.Flags().StringArrayVar(&ExecFlag.GidMappings, "gidmapping", []string{}, "source:target:[flags:[fstype:[option]]]")
	cmd.Flags().IntVar(&ExecFlag.UID, "uid", os.Getuid(), "用户ID")
	cmd.Flags().IntVar(&ExecFlag.GID, "gid", os.Getuid(), "用户组ID")
	cmd.Flags().StringVar(&ExecFlag.Socket, "socket", "", "可重入终端通信套接字,指定相同的套接字将重用已启动的环境")
	cmd.Flags().StringVar(&ExecFlag.RootFS, "rootfs", "", "合并的根目录位置")
	cmd.Flags().BoolVar(&ExecFlag.NoBindRootFS, "no-bind-rootfs", false, "手动挂载rootfs")
	cmd.Flags().BoolVar(&ExecFlag.AutoExit, "auto-exit", true, "当没有进程连接时，自动退出服务")
	cmd.Flags().BoolVar(&ExecFlag.NoFail, "no-fail", false, "任何步骤失败时立即退出")
	cmd.Flags().BoolVar(&ExecFlag.Root, "root", false, "以root身份运行（覆盖uid/gid选项）")
	cmd.Flags().DurationVar(&ExecFlag.SocketTimeout, "socket-timeout", 30*time.Second, "终端套接字连接超时")
	cmd.Flags().StringVar(&ExecFlag.FuseOverlayFS, "fuse-overlayfs", "", "外部fuse-overlayfs命令路径(可选)")
	cmd.Flags().StringVar(&ExecFlag.FuseOverlayFSArgs, "fuse-overlayfs-args", "", "fuse-overlayfs命令额外参数")
	cmd.Flags().SortFlags = false
	return cmd
}
