/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _export

import (
	"ll-killer/utils"

	"github.com/spf13/cobra"
)

var ExportFlag struct {
	Self  string
	Shell string
	Args  []string
}

func ExportMain(cmd *cobra.Command, args []string) error {
	args = append([]string{"ll-builder", "export"}, args...)
	utils.Exec(args...)
	return nil
}

func CreateExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "export",
		Short:              "导出应用",
		Long:               "此命令执行ll-builder export，用于提供一致性体验。",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(ExportMain(cmd, args))
		},
	}

	return cmd
}
