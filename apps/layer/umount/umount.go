/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _umount

import (
	"github.com/System233/ll-killer-go/layer"
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

var Flag layer.UmountOption

func UmountMain(cmd *cobra.Command, args []string) error {
	return layer.Umount(&Flag)
}

func CreateUmountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "umount <layer挂载目录> [flags] -- [fusermount选项]",
		Short: "卸载layer挂载点。",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Flag.Target = args[0]
			Flag.Args = args[1:]
			utils.ExitWith(UmountMain(cmd, args[1:]))
		},
	}
	cmd.Flags().StringVar(&Flag.ExecPath, "exec", layer.FuserMount, "指定fusermount命令位置")
	return cmd
}
