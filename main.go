/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package main

import (
	"fmt"
	_apt "ll-killer/apps/apt"
	_build "ll-killer/apps/build"
	_buildaux "ll-killer/apps/build-aux"
	_clean "ll-killer/apps/clean"
	_commit "ll-killer/apps/commit"
	_create "ll-killer/apps/create"
	_exec "ll-killer/apps/exec"
	_export "ll-killer/apps/export"
	_layer "ll-killer/apps/layer"
	_nsenter "ll-killer/apps/nsenter"
	_overlay "ll-killer/apps/overlay"
	_ptrace "ll-killer/apps/ptrace"
	_run "ll-killer/apps/run"
	_script "ll-killer/apps/script"
	"ll-killer/config"
	"ll-killer/utils"
	"log"
	"os"

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

运行 "ll-killer <command> --help" 以查看子命令的详细信息。

更多信息请查看项目主页: https://github.com/System233/ll-killer-go.git
`

func main() {
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
		Use:     "ll-killer",
		Short:   Usage,
		Example: utils.BuildHelpMessage(MainCommandHelp),
	}
	app.Flags().SortFlags = false
	app.InheritedFlags().SortFlags = false
	app.LocalFlags().SortFlags = false
	app.Flags().BoolVar(&utils.GlobalFlag.Debug, "debug", utils.GlobalFlag.Debug, "显示调试信息")
	app.AddCommand(_apt.CreateAPTCommand(),
		_build.CreateBuildCommand(),
		_exec.CreateExecCommand(),
		_run.CreateRunCommand(),
		_create.CreateCreateCommand(),
		_commit.CreateCommitCommand(),
		_layer.CreateLayerCommand(),
		_clean.CreateCleanCommand(),
		_export.CreateExportCommand(),
		_buildaux.CreateBuildAuxCommand(),
		_script.CreateScriptCommand(),
		_overlay.CreateOverlayCommand(),
		_nsenter.NsEnterNsEnterCommand())
	app.Version = fmt.Sprintf("%s/%s", config.Version, config.BuildTime)
	if _ptrace.IsSupported {
		app.AddCommand(_ptrace.CreatePtraceCommand())
	}
	if err := app.Execute(); err != nil {
		utils.ExitWith(err)
	}

}
