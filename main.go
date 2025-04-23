/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package main

import (
	"fmt"
	"log"
	"os"

	_apt "github.com/System233/ll-killer-go/apps/apt"
	_build "github.com/System233/ll-killer-go/apps/build"
	_buildaux "github.com/System233/ll-killer-go/apps/build-aux"
	_clean "github.com/System233/ll-killer-go/apps/clean"
	_commit "github.com/System233/ll-killer-go/apps/commit"
	_create "github.com/System233/ll-killer-go/apps/create"
	_exec "github.com/System233/ll-killer-go/apps/exec"
	_export "github.com/System233/ll-killer-go/apps/export"
	_layer "github.com/System233/ll-killer-go/apps/layer"
	_nsenter "github.com/System233/ll-killer-go/apps/nsenter"
	_overlay "github.com/System233/ll-killer-go/apps/overlay"
	_ptrace "github.com/System233/ll-killer-go/apps/ptrace"
	_run "github.com/System233/ll-killer-go/apps/run"
	_script "github.com/System233/ll-killer-go/apps/script"
	_update "github.com/System233/ll-killer-go/apps/update"
	"github.com/System233/ll-killer-go/config"
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

const (
	Usage = `ll-killer - 玲珑容器辅助工具`
)
const MainCommandHelp = `ll-killer 是一个工具，旨在解决玲珑容器应用的构建问题。

项目构建一般经历以下几个过程：
  create  创建项目，生成必要的构建文件。
  build   进入构建环境，执行构建操作，如apt安装等。
  commit  提交构建内容至玲珑容器。
  run     运行已构建的应用进行测试。

* 关于如何加速构建，优化磁盘占用，请查看layer子命令帮助。

运行 "<program> <command> --help" 以查看子命令的详细信息。

更多信息请查看项目主页: https://github.com/System233/ll-killer-go.git
`

func main() {
	if err := utils.InitChildSubreaper(); err != nil {
		utils.ExitWith(err)
	}
	if os.Getenv(config.KillerDebug) != "" {
		utils.GlobalFlag.Debug = true
	}
	err := utils.SetupEnvVar()
	if err != nil {
		utils.Debug("SetupEnvVar:", err)
	}
	cobra.EnableCommandSorting = false
	log.SetFlags(0)

	if utils.GlobalFlag.Debug {
		pid := os.Getpid()
		log.SetPrefix(fmt.Sprintf("[PID %d] ", pid))
	}

	app := cobra.Command{
		Use:     config.KillerExec,
		Short:   Usage,
		Example: utils.BuildHelpMessage(MainCommandHelp),
	}
	app.Flags().SortFlags = false
	app.InheritedFlags().SortFlags = false
	app.LocalFlags().SortFlags = false
	app.PersistentFlags().BoolVar(&utils.GlobalFlag.Debug, "debug", utils.GlobalFlag.Debug, "显示调试信息")
	app.AddCommand(
		_buildaux.CreateBuildAuxCommand(),
		_create.CreateCreateCommand(),
		_apt.CreateAPTCommand(),
		_build.CreateBuildCommand(),
		_exec.CreateExecCommand(),
		_layer.CreateLayerCommand(),
		_clean.CreateCleanCommand(),
		_script.CreateScriptCommand(),
		_overlay.CreateOverlayCommand(),
		_nsenter.NsEnterNsEnterCommand(),
		_run.CreateRunCommand(),
		_commit.CreateCommitCommand(),
		_export.CreateExportCommand(),
		_update.CreateUpdateCommand())
	app.Version = fmt.Sprintf("%s/%s", config.Version, config.BuildTime)
	if _ptrace.IsSupported {
		app.AddCommand(_ptrace.CreatePtraceCommand())
	}
	utils.ExitWith(app.Execute())

}
