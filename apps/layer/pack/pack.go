/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _pack

import (
	"os"

	"github.com/System233/ll-killer-go/layer"
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

const PackCommandHelp = `将指定的layer结构文件夹打包为layer,此命令使用mkfs.erofs对layer结构目录进行打包，从而无需ll-builder提交到ostree，减少不必要的文件复制。

本功能的目的是消除~/.cache/linglong-builder下不必要的应用缓存，保护磁盘寿命。
* 源文件夹需要为合法的layer结构，参考linglong/output/binary目录。

目前实际可操作的流程如下：

* 使用ll-builder build --skip-commit-output获取原始输出linglong/output/_build/，
  并手动创建类似linglong/output/binary的目录结构，最后使用本命令打包。
  注意：原始输出未对快捷方式等添加'll-cli run {APPID} -- '前缀。
* 挂载已有的layer，复制layer内容并修改，最后重新打包。
* 挂载已有的layer，使用overlayfs叠加layer目录，在合并目录中修改，最后对合并目录重新打包。
* 挂载已有的高版本layer，使用lz4hc算法重新打包，避免旧版玲珑因erofs不支持zstd压缩算法而无法安装的问题。
* 使用'<program> layer build'子命令自动调用本功能。

`

var Flag layer.PackOption

func PackMain(cmd *cobra.Command, args []string) error {
	utils.Debug("PackMain", args)
	Flag.Args = args
	return layer.Pack(&Flag)
}

func CreatePackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pack <layer结构文件夹> [flags] -- [mkfs.erofs选项]",
		Short:         "将指定的layer结构文件夹打包为layer。",
		Long:          utils.BuildHelpMessage(PackCommandHelp),
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Flag.Source = args[0]
			utils.ExitWith(PackMain(cmd, args[1:]))
		},
	}
	cmd.Flags().IntVarP(&Flag.BlockSize, "block-size", "b", 4096, "块大小")
	cmd.Flags().StringVarP(&Flag.Target, "output", "o", "", "输出文件名,默认与玲珑命名方式相同。")
	cmd.Flags().StringVar(&Flag.ExecPath, "exec", layer.MkfsErofs, "指定mkfs.erofs命令位置")
	cmd.Flags().IntVarP(&Flag.Uid, "force-uid", "U", os.Getuid(), "文件Uid,-1为不更改")
	cmd.Flags().IntVarP(&Flag.Gid, "force-gid", "G", os.Getegid(), "文件Gid,-1为不更改")
	cmd.Flags().StringVarP(&Flag.Compressor, "compressor", "z", "lz4hc", "压缩算法，请查看mkfs.erofs帮助")
	cmd.Flags().SortFlags = false
	return cmd
}
