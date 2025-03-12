/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _buildaux

import (
	"embed"
	"errors"
	"io/fs"
	"ll-killer/config"
	"ll-killer/utils"
	"log"
	"os"
	"path/filepath"

	"github.com/moby/sys/reexec"
	"github.com/spf13/cobra"
)

//go:embed build-aux apt.conf.d sources.list.d/.keep
var content embed.FS

const BuildAuxCommandHelp = `
ll-killer build-aux创建一系列辅助脚本，可用于构建和调试：

build-aux 目录下创建的工具：
  - entrypoint.sh        玲珑应用入口点
  - env.sh               运行环境变量配置
  - ldd-check.sh         检查容器内缺失库（处理未完整声明依赖的 deb）
  - ldd-search.sh        在 ll-killer apt 环境中搜索缺失库所属 deb 包
  - relink.sh            修复不支持相对路径的符号链接
  - setup-desktop.sh     修复 .desktop 文件的 Icon 和 Exec 路径
  - setup-filesystem.sh  从构建环境复制文件到 $PREFIX
  - setup-icon.sh        处理图标文件，支持 ico/png/jpg/gif/svg 格式
  - setup.sh             执行所有修复操作并完成文件复制

`

var BuildAuxFlag struct {
	Force bool
}

func embedFilesToDisk(destDir string, force bool) error {
	err := fs.WalkDir(content, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, path)

		if !d.IsDir() {
			if !force && utils.IsExist(destPath) {
				log.Println("skip:", destPath)
				return nil
			}
			srcFile, err := content.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			err = utils.CopyFile(destPath, srcFile, 0755, force)
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
	if err != nil {
		return err
	}
	err = utils.CopySymlink("Makefile", "build-aux/Makefile", force)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Println("skip:", "Makefile")
			return nil
		}
		return err
	}
	log.Println("created:", "Makefile")
	return nil
}

func SetupKillerExec(target string) error {
	if !BuildAuxFlag.Force && utils.IsExist(config.KillerExec) {
		log.Println("skip:", target)
		return nil
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	isSame, err := utils.IsSameFile(target, self)
	if err != nil {
		return err
	}
	if isSame {
		log.Println("skip same:", target)
		return nil
	}
	err = utils.CopyFileIO(reexec.Self(), target)
	if err != nil {
		return err
	}
	log.Println("created:", target)
	return nil
}
func ExtractBuildAuxFiles(force bool) error {
	if err := embedFilesToDisk(".", force); err != nil {
		return err
	}
	if err := SetupKillerExec(config.KillerExec); err != nil {
		return err
	}
	return nil
}

func BuildAuxMain(cmd *cobra.Command, args []string) error {
	return ExtractBuildAuxFiles(BuildAuxFlag.Force)
}

func CreateBuildAuxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build-aux",
		Short: "创建辅助构建脚本",
		Long:  BuildAuxCommandHelp,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(BuildAuxMain(cmd, args))
		},
	}
	cmd.Flags().BoolVar(&BuildAuxFlag.Force, "force", false, "强制覆盖已存在文件")
	return cmd
}
