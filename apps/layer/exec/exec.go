/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _exec

import (
	"fmt"
	"syscall"
	"time"

	"github.com/System233/ll-killer-go/apps/layer/fs"
	"github.com/System233/ll-killer-go/layer"
	"github.com/System233/ll-killer-go/reexec"
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

var Flag struct {
	RootFs      string
	Runtime     string
	Quiet       bool
	Wait        bool
	WaitTimeout time.Duration
	Args        []string
}

const ExecCommandDescription = `此命令用于提供一个指定base和runtime的玲珑容器模拟环境，以实现更灵活高效的应用构建。
* 当build子命令指定了rootfs参数时会自动调用本命令。
* 与layer build子命令相似，但本命令专注于启动一个容器环境。
`
const ExecCommandHelp = ``

var Config layer.Config
var LayerInfo layer.LayerInfo

func GetExecArgs() []string {
	args := []string{
		fmt.Sprint("--quiet=", Flag.Quiet),
	}
	if Flag.RootFs != "" {
		args = append(args, "--rootfs", Flag.RootFs)
	}
	if Flag.Runtime != "" {
		args = append(args, "--runtime", Flag.Runtime)
	}
	args = append(args, fmt.Sprint("--wait=", Flag.Wait))

	if Flag.WaitTimeout != 0 {
		args = append(args, "--wait-timeout", fmt.Sprint(Flag.WaitTimeout))
	}
	if len(Flag.Args) > 0 {
		args = append(args, "--")
		args = append(args, Flag.Args...)
	}
	return args
}
func ExecLayer() error {
	if err := utils.RemountProc(); err != nil {
		return err
	}

	return fs.Run(fs.SetupFilesystemOption{
		RootFs:    Flag.RootFs,
		Runtime:   Flag.Runtime,
		Quiet:     Flag.Quiet,
		Config:    &Config,
		LayerInfo: &LayerInfo,
	}, Flag.Args...)
}
func ExecMain(cmd *cobra.Command, args []string) error {
	Flag.Args = args

	if Flag.Wait {
		utils.SetChildSubreaperWaitDuration(Flag.WaitTimeout)
	}

	reexec.Register("ExecLayer", ExecLayer)
	ok, err := reexec.Init()
	if ok || err != nil {
		return err
	}
	return utils.SwitchTo("ExecLayer", &utils.SwitchFlags{
		UID:           0,
		GID:           0,
		Cloneflags:    syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID,
		Args:          append([]string{"layer", "exec"}, GetExecArgs()...),
		NoDefaultArgs: true,
	})
}

func CreateExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exec [flags] -- cmd",
		Short:   "进入自定义的构建环境，相当于不导出layer的build。",
		Long:    utils.BuildHelpMessage(ExecCommandDescription),
		Example: utils.BuildHelpMessage(ExecCommandHelp),
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(ExecMain(cmd, args))
		},
	}
	cmd.Flags().StringVar(&Flag.RootFs, "rootfs", "/", "根文件系统")
	cmd.Flags().StringVar(&Flag.Runtime, "runtime", "", "runtime文件系统")
	cmd.Flags().BoolVar(&Flag.Wait, "wait", false, "作为服务进程等待所有进程退出，默认杀死所有子进程。")
	cmd.Flags().DurationVar(&Flag.WaitTimeout, "wait-timeout", -1, "等待所有进程退出的最大超时时间，-1为永久等待。")
	cmd.Flags().BoolVar(&Flag.Quiet, "quiet", true, "安静模式，构建前不输出项目信息")
	cmd.Flags().SortFlags = false
	return cmd
}
