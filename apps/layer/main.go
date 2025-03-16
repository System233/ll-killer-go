/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _layer

import (
	_build "github.com/System233/ll-killer-go/apps/layer/build"
	_dump "github.com/System233/ll-killer-go/apps/layer/dump"
	_exec "github.com/System233/ll-killer-go/apps/layer/exec"
	_mount "github.com/System233/ll-killer-go/apps/layer/mount"
	_pack "github.com/System233/ll-killer-go/apps/layer/pack"
	_umount "github.com/System233/ll-killer-go/apps/layer/umount"

	"github.com/spf13/cobra"
)

func CreateLayerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "layer",
		Short: "layer 构建/打包/挂载/调试工具",
		Long:  "独立于玲珑的layer管理器，提供强大的layer文件处理支持，需要安装erofs-utils。",
	}
	cmd.AddCommand(
		_pack.CreatePackCommand(),
		_build.CreateBuildCommand(),
		_exec.CreateExecCommand(),
		_mount.CreateMountCommand(),
		_umount.CreateUmountCommand(),
		_dump.CreateDumpCommand())
	return cmd
}
