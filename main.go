/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package main

import (
	"fmt"
	"ll-killer/ptrace"
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
	if os.Getenv(kKillerDebug) != "" {
		GlobalFlag.Debug = true
	}
	err := SetupEnvVar()
	if err != nil {
		Debug("SetupEnvVar:", err)
	}
	cobra.EnableCommandSorting = false
	// log.SetFlags(0)
	app := cobra.Command{
		Use:     "ll-killer",
		Short:   Usage,
		Example: BuildHelpMessage(MainCommandHelp),
	}
	app.Flags().SortFlags = false
	app.InheritedFlags().SortFlags = false
	app.LocalFlags().SortFlags = false
	app.Flags().BoolVar(&GlobalFlag.Debug, "debug", GlobalFlag.Debug, "显示调试信息")
	app.AddCommand(CreateAPTCommand(),
		CreateBuildCommand(),
		CreateExecCommand(),
		CreateRunCommand(),
		CreateCreateCommand(),
		CreateCommitCommand(),
		CreateCleanCommand(),
		CreateExportCommand(),
		CreateBuildAuxCommand(),
		CreateScriptCommand(),
		CreateOverlayCommand())
	app.Version = fmt.Sprintf("%s/%s", Version, BuildTime)
	if ptrace.IsSupported {
		app.AddCommand(CreatePtraceCommand())
	}
	if err := app.Execute(); err != nil {
		ExitWith(err)
	}

}
