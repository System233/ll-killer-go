package fs

import (
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/System233/ll-killer-go/config"
	"github.com/System233/ll-killer-go/layer"
	"github.com/System233/ll-killer-go/utils"
	"golang.org/x/sys/unix"
)

const WorkDir = "linglong/output"

type SetupFilesystemOption struct {
	RootFs    string
	Runtime   string
	Config    *layer.Config
	LayerInfo *layer.LayerInfo
	Quiet     bool
}

func SetupFilesystem(opt SetupFilesystemOption) error {
	err := utils.LoadYamlFile(config.LinglongYaml, opt.Config)
	if err != nil {
		return err
	}
	if err := opt.LayerInfo.ParseLayerInfo(*opt.Config); err != nil {
		return fmt.Errorf("解析yaml错误:%v", err)
	}
	if !opt.Quiet {
		opt.LayerInfo.Print()
	}
	if utils.IsExist(WorkDir) {
		if err := os.RemoveAll(WorkDir); err != nil {
			return fmt.Errorf("无法移除%s:%v", WorkDir, err)
		}
	}
	if err := os.MkdirAll(WorkDir, 0755); err != nil {
		return fmt.Errorf("无法建立工作目录%s:%v", WorkDir, err)
	}

	configID := opt.Config.Package.ID
	configBuild := opt.Config.Build
	rootfsPath := path.Join(WorkDir, "rootfs")

	if err := utils.Mount(&utils.MountOption{Source: opt.RootFs, Target: rootfsPath, FSType: "merge", Flags: unix.MS_RDONLY}); err != nil {
		return fmt.Errorf("挂载根目录失败:%v", err)
	}
	if opt.RootFs != "/" {
		if err := utils.MountAll([]utils.MountOption{
			{
				Source: "/dev",
				Target: path.Join(rootfsPath, "dev"),
			},
			{
				Source: "/proc",
				Target: path.Join(rootfsPath, "proc"),
			},
			{
				Source: "/home",
				Target: path.Join(rootfsPath, "home"),
			},
			{
				Source: "/root",
				Target: path.Join(rootfsPath, "root"),
			},
			{
				Source: "/tmp",
				Target: path.Join(rootfsPath, "tmp"),
			},
			{
				Source: "/sys",
				Target: path.Join(rootfsPath, "sys"),
			}}); err != nil {
			return fmt.Errorf("挂载主机文件系统失败:%v", err)
		}
	}

	if err := utils.MountAll([]utils.MountOption{
		{
			Source: "tmpfs",
			Target: path.Join(rootfsPath, "run"),
			FSType: "tmpfs",
		},
		{
			Source: "/run/systemd",
			Target: path.Join(rootfsPath, "run/systemd"),
		},
		{
			Source: "/etc/resolv.conf",
			Target: path.Join(rootfsPath, "etc/resolv.conf"),
		},
		{
			Source: "/etc/localtime",
			Target: path.Join(rootfsPath, "etc/localtime"),
		},
		{
			Source: "/etc/timezone",
			Target: path.Join(rootfsPath, "etc/timezone"),
		},
		{
			Source: "/etc/machine-id",
			Target: path.Join(rootfsPath, "etc/machine-id"),
		},
	}); err != nil {
		return fmt.Errorf("挂载主机配置文件失败:%v", err)
	}

	if err := utils.MountAll([]utils.MountOption{
		{
			Source: ".",
			Target: path.Join(rootfsPath, "project"),
		},
		{
			Source: config.AptDataDir,
			Target: path.Join(rootfsPath, "var/lib/apt"),
		},
		{
			Source: config.AptCacheDir,
			Target: path.Join(rootfsPath, "var/cache"),
		}}); err != nil {
		return fmt.Errorf("挂载项目目录失败:%v", err)
	}

	if opt.Runtime != "" {
		runtimeFS := path.Join(rootfsPath, "runtime")
		if err := utils.Mount(&utils.MountOption{Source: opt.Runtime, Target: runtimeFS}); err != nil {
			return fmt.Errorf("挂载runtime目录失败:%v", err)
		}
	}

	runHostRootfs := path.Join(rootfsPath, "run/host/rootfs")
	if err := os.MkdirAll(runHostRootfs, 0755); err != nil {
		return fmt.Errorf("创建run/host/rootfs目录失败:%v", err)
	}

	buildHostDir := path.Join(WorkDir, "build")
	if err := os.MkdirAll(buildHostDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败:%v", err)
	}

	optTmpfs := path.Join(rootfsPath, "opt")
	if err := utils.Mount(&utils.MountOption{Source: "tmpfs", Target: optTmpfs, FSType: "tmpfs"}); err != nil {
		return fmt.Errorf("挂载opt目录失败:%v", err)
	}
	optAppsDir := path.Join(rootfsPath, "opt/apps", configID, "files")
	if err := os.MkdirAll(optAppsDir, 0755); err != nil {
		return fmt.Errorf("创建/opt/apps目录失败:%v", err)
	}
	if err := utils.MountBind(buildHostDir, optAppsDir, syscall.MS_BIND); err != nil {
		return fmt.Errorf("挂载输出目录失败:%v", err)
	}

	entryPath := "linglong/entry.sh"
	entryData := []byte(fmt.Sprintf("#!/bin/bash\n%s", configBuild))
	if err := os.WriteFile(entryPath, entryData, 0755); err != nil {
		return fmt.Errorf("写入entry.sh失败:%v", err)
	}

	if err := utils.MountBind(rootfsPath, rootfsPath, 0); err != nil {
		return fmt.Errorf("绑定根目录失败:%v", err)
	}

	if err := syscall.PivotRoot(rootfsPath, runHostRootfs); err != nil {
		return fmt.Errorf("切换根目录失败:%v", err)
	}

	os.Setenv("LINGLONG_APPID", configID)
	os.Setenv("PREFIX", path.Join("/opt/apps", configID, "files"))
	os.Setenv("TRIPLET", layer.GetTriplet())
	os.Setenv(config.KillerPackerEnv, "1")
	return nil
}

func PostFilesystem() error {
	runHostRootfs := "/run/host/rootfs"
	rootfsPath := path.Join(WorkDir, "rootfs")
	if err := unix.PivotRoot(runHostRootfs, rootfsPath); err != nil {
		return fmt.Errorf("切换回主机失败:%v", err)
	}
	return nil
}
func Run(opt SetupFilesystemOption, args ...string) error {
	err := SetupFilesystem(opt)
	if err != nil {
		return err
	}
	defer PostFilesystem()
	return utils.ExecRaw(args...)
}
