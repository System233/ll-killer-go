/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _overlay

import (
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

func OverlayMain(cmd *cobra.Command, args []string) error {
	return utils.FuseOvlMain(args)
}

func CreateOverlayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "overlay",
		Short:              "内置fuse-overlayfs挂载",
		Long:               "此命令调用内嵌的fuse-overlayfs程序实现overlay挂载",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(OverlayMain(cmd, args))
		},
	}

	return cmd
}
