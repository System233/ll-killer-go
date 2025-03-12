/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _script

import (
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

var ScriptFlag struct {
}

const ScriptCommandHelp = `
此命令为命令提供KILLER_EXEC环境变量，指向ll-killer二进制的绝对路径，确保构建脚本能够找到ll-killer。
用法：
<program> script -- <构建脚本> [参数...]

此命令等同于: 
KILLER_EXEC=<program> <构建脚本> [参数...]
`

func ScriptMain(cmd *cobra.Command, args []string) error {
	return utils.ExecRaw(args...)
}

func CreateScriptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "script",
		Short: "执行自定义构建流程",
		Long:  utils.BuildHelpMessage(ScriptCommandHelp),
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(ScriptMain(cmd, args))
		},
	}

	return cmd
}
