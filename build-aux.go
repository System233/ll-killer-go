/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed build-aux apt.conf.d sources.list.d/.keep
var content embed.FS

func embedFilesToDisk(destDir string) error {
	err := fs.WalkDir(content, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, path)

		if !d.IsDir() {
			if IsExist(destPath) {
				log.Println("skip:", destPath)
				return nil
			}
			srcFile, err := content.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}

			log.Println("created:", destPath)
		} else {
			err = os.MkdirAll(destPath, 0755)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func BuildAuxMain(cmd *cobra.Command, args []string) error {
	return embedFilesToDisk(".")
}

func CreateBuildAuxCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "build-aux",
		Short: "创建辅助构建脚本",
		Run: func(cmd *cobra.Command, args []string) {
			ExitWith(BuildAuxMain(cmd, args))
		},
	}
}
