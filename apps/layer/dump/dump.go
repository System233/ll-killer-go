/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _dump

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/System233/ll-killer-go/layer"
	"github.com/System233/ll-killer-go/utils"
	"github.com/spf13/cobra"
)

var Flag struct {
	Target     string
	ShowAll    bool
	ShowHeader bool
	ShowLayer  bool
	ShowErofs  bool
	ExecPath   string
	Format     string
	ShowJson   bool
}

func DumpMain(cmd *cobra.Command, args []string) error {
	header, err := layer.NewLayerHeaderFromFile(Flag.Target)
	if err != nil {
		return err
	}
	if Flag.Format != "" {
		t, err := template.New("format").Parse(Flag.Format)
		if err != nil {
			return fmt.Errorf("解析format失败:%v", err)
		}
		return t.Execute(os.Stdout, header)
	}
	if Flag.ShowJson {
		return json.NewEncoder(os.Stdout).Encode(header)
	}
	if Flag.ShowHeader {
		header.Print()
	}
	if Flag.ShowLayer {
		header.Info.Print()
	}
	if Flag.ShowErofs {
		err = header.PrintErofs(&layer.DumpErofsOption{
			ExecPath: Flag.ExecPath,
			Args:     args,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump <layer文件> [flags] -- [dump.erofs选项]",
		Short: "输出layer信息。",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Flag.Target = args[0]
			if Flag.ShowErofs || Flag.ShowHeader || Flag.ShowLayer {
				Flag.ShowAll = false
			}
			if Flag.ShowAll {
				Flag.ShowHeader = true
				Flag.ShowLayer = true
				Flag.ShowErofs = true
			}
			utils.ExitWith(DumpMain(cmd, args[1:]))
		},
	}
	cmd.Flags().BoolVarP(&Flag.ShowHeader, "header", "x", false, "显示文件头信息")
	cmd.Flags().BoolVarP(&Flag.ShowLayer, "layer", "l", false, "显示Layer信息")
	cmd.Flags().BoolVarP(&Flag.ShowErofs, "erofs", "e", false, "显示Erofs信息")
	cmd.Flags().BoolVarP(&Flag.ShowAll, "all", "a", true, "显示全部信息")
	cmd.Flags().BoolVarP(&Flag.ShowJson, "json", "j", false, "输出json")
	cmd.Flags().StringVarP(&Flag.Format, "format", "f", "", "格式化输出，字段详见json")
	cmd.Flags().StringVar(&Flag.ExecPath, "exec", layer.DumpErofs, "指定dump.erofs命令的路径")
	cmd.Flags().SortFlags = false
	return cmd
}
