/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _buildaux

import (
	buildaux "github.com/System233/ll-killer-go/build-aux"
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

const BuildAuxCommandHelp = `
ll-killer build-aux创建一系列辅助脚本，可用于构建和调试：

build-aux 目录下创建的工具：
  - entrypoint.sh        玲珑应用入口点
  - env.sh               运行环境变量配置
  - ldd-check.sh         检查容器内缺失库（处理未完整声明依赖的 deb）
  - ldd-search.sh        在 ll-killer apt 环境中搜索缺失库所属 deb 包
  - relink.sh            修复不支持相对路径的符号链接
  - setup-desktop.sh     修复 .desktop 文件的 Icon 和 Exec 路径
  - setup-filesystem.sh  从构建环境复制文件到 $PREFIX
  - setup-icon.sh        处理图标文件，支持 ico/png/jpg/gif/svg 格式
  - setup.sh             执行所有修复操作并完成文件复制

`

var BuildAuxFlag struct {
	Force bool
}

func BuildAuxMain(cmd *cobra.Command, args []string) error {
	return buildaux.ExtractBuildAuxFiles(BuildAuxFlag.Force)
}

func CreateBuildAuxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build-aux",
		Short: "创建辅助构建脚本",
		Long:  BuildAuxCommandHelp,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(BuildAuxMain(cmd, args))
		},
	}
	cmd.Flags().BoolVar(&BuildAuxFlag.Force, "force", false, "强制覆盖已存在文件")
	return cmd
}
