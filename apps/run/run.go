/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _run

import (
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

var RunFlag struct {
	Self  string
	Shell string
	Args  []string
}

func RunMain(cmd *cobra.Command, args []string) error {
	args = append([]string{"ll-builder", "run"}, args...)
	utils.Exec(args...)
	return nil
}

func CreateRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "run",
		Short:              "启动容器",
		Long:               "此命令执行ll-builder run，用于提供一致性体验。",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(RunMain(cmd, args))
		},
	}

	return cmd
}
