/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _mount

import (
	"ll-killer/layer"
	"ll-killer/utils"

	"github.com/spf13/cobra"
)

const MountCommandHelp = `
`

var Flag layer.MountOption

func MountMain(cmd *cobra.Command, args []string) error {
	Flag.Args = args
	return layer.Mount(&Flag)
}

func CreateMountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mount <layer文件名> <源目录>  -- [erofsfuse选项]",
		Short: "挂载layer文件。",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			i := 1
			Flag.Source = args[0]
			if len(args) > 1 {
				Flag.Target = args[1]
				i++
			}
			utils.ExitWith(MountMain(cmd, args[i:]))
		},
	}
	cmd.Flags().StringVar(&Flag.ExecPath, "exec", layer.ErofsFuse, "指定erofsfuse命令位置")
	return cmd
}
