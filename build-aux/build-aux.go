/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package buildaux

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
)

//go:embed build-aux apt.conf.d sources.list.d/.keep
var content embed.FS

func ExtractEmbedFilesToDisk(destDir string, force bool) error {
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

func ExtractKillerExec(target string, force bool) error {
	if !force && utils.IsExist(config.KillerExec) {
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
	if err := ExtractEmbedFilesToDisk(".", force); err != nil {
		return err
	}
	if err := ExtractKillerExec(config.KillerExec, force); err != nil {
		return err
	}
	return nil
}
