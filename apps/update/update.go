/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _update

import (
	"time"

	"github.com/System233/ll-killer-go/updater"
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

var Flag updater.UpdateOption

const ScriptCommandHelp = ``

func UpdateMain(cmd *cobra.Command, args []string) error {
	return updater.Update(Flag)
}

func CreateUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "更新ll-killer",
		Long:  utils.BuildHelpMessage(ScriptCommandHelp),
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(UpdateMain(cmd, args))
		},
	}
	cmd.Flags().BoolVarP(&Flag.Yes, "yes", "y", false, "直接执行更新，不询问是和否。")
	cmd.Flags().IntVar(&Flag.Retry, "retry", 10, "最大重试次数")
	cmd.Flags().DurationVar(&Flag.Timeout, "timeout", time.Second*5, "最大重试次数")
	cmd.Flags().SortFlags = false
	return cmd
}
