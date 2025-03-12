/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _build

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	_ptrace "github.com/System233/ll-killer-go/apps/ptrace"
	"github.com/System233/ll-killer-go/config"
	"github.com/System233/ll-killer-go/utils"

	"github.com/moby/sys/reexec"
	"github.com/spf13/cobra"
)

var BuildFlag struct {
	RootFS            string
	TmpRootFS         string
	CWD               string
	Strict            bool
	Ptrace            bool
	EncodedArgs       string
	Self              string
	Args              []string
	FuseOverlayFS     string
	FuseOverlayFSArgs string
}

const BuildCommandDescription = `进入玲珑构建环境，可执行 apt 安装等构建操作。`
const BuildCommandHelp = `
# 直接运行 ll-killer build 进入交互式构建环境
<program> build

# 使用 ll-killer build -- <命令> 直接执行指定构建命令
<program> build

# 使用 fuse-overlayfs 进行构建，指定 fuse-overlayfs 命令路径
<program> build --fuse-overlayfs /path/to/fuse-overlayfs

# 使用 fuse-overlayfs，并传入额外的命令参数
<program> build --fuse-overlayfs-args "--option=value"

# 启用严格模式，确保构建环境与运行时环境一致，且不包含开发工具（默认行为）
<program> build --strict

`

func GetBuildArgs() []string {
	args := []string{}

	if BuildFlag.RootFS != "" {
		args = append(args, "--rootfs", BuildFlag.RootFS)
	}

	if BuildFlag.CWD != "" {
		args = append(args, "--cwd", BuildFlag.CWD)
	}

	if BuildFlag.Strict {
		args = append(args, "--strict")
	}
	if _ptrace.IsSupported {
		args = append(args, "--ptrace="+fmt.Sprint(BuildFlag.Ptrace))
	}

	if BuildFlag.FuseOverlayFS != "" {
		args = append(args, "--fuse-overlayfs", BuildFlag.FuseOverlayFS)
	}

	if BuildFlag.FuseOverlayFSArgs != "" {
		args = append(args, "--fuse-overlayfs-args", BuildFlag.FuseOverlayFSArgs)
	}

	if BuildFlag.Self != "" {
		args = append(args, "--self", BuildFlag.Self)
	}

	if len(BuildFlag.Args) > 0 {
		args = append(args, "--")
		args = append(args, BuildFlag.Args...)
	}

	return args
}

func MountOverlayStage2() {

	overlayDir := path.Join(BuildFlag.CWD, config.FileSystemDir)
	mergedDir := path.Join(overlayDir, "merged")
	tmpRootFS := BuildFlag.TmpRootFS
	err := syscall.PivotRoot(tmpRootFS+mergedDir, tmpRootFS+mergedDir+BuildFlag.RootFS)
	if err != nil {
		utils.ExitWith(err, "PivotRoot2")
	}
	if BuildFlag.Ptrace && _ptrace.IsSupported {
		_ptrace.Ptrace(BuildFlag.Self, BuildFlag.Args)
	} else {
		utils.Exec(BuildFlag.Args...)
	}

}
func MountOverlay() {
	overlayDir := path.Join(BuildFlag.CWD, config.FileSystemDir)
	aptCacheDir := path.Join(BuildFlag.CWD, config.AptCacheDir)
	aptDataDir := path.Join(BuildFlag.CWD, config.AptDataDir)
	tmpRootFS := BuildFlag.TmpRootFS
	upperDir := path.Join(overlayDir, config.UpperDirName)
	lowerDir := path.Join(overlayDir, config.LowerDirName)
	workDir := path.Join(overlayDir, config.WorkDirName)
	mergedDir := path.Join(overlayDir, config.MergedDirName)
	cwdRootFSPivoted := fmt.Sprint(BuildFlag.RootFS, tmpRootFS)
	err := utils.MkdirAlls([]string{
		tmpRootFS, upperDir, lowerDir, workDir,
		mergedDir,
		aptCacheDir,
		aptDataDir,
	}, 0755)
	if err != nil {
		utils.ExitWith(err)
	}
	err = utils.MountBind(BuildFlag.RootFS, BuildFlag.RootFS, 0)
	if err != nil {
		utils.ExitWith(err)
	}

	err = utils.MountBind("/", tmpRootFS, 0)
	if err != nil {
		utils.ExitWith(err)
	}

	err = utils.MountAll([]utils.MountOption{
		{
			Source: "sources.list.d",
			Target: lowerDir + "/etc/apt/sources.list.d",
		},
		{
			Source: "sources.list",
			Target: lowerDir + "/etc/apt/sources.list",
		},
		{
			Source: "apt.conf",
			Target: lowerDir + "/etc/apt/apt.conf",
		},
		{
			Source: "apt.conf.d",
			Target: lowerDir + "/etc/apt/apt.conf.d",
		},
		{
			Source: "auth.conf.d",
			Target: lowerDir + "/etc/apt/auth.conf.d",
		},
	})
	if err != nil {
		utils.ExitWith(err, "MountAll")
	}
	err = syscall.PivotRoot(BuildFlag.RootFS, cwdRootFSPivoted)
	if err != nil {
		utils.ExitWith(err, "PivotRoot", BuildFlag.RootFS, cwdRootFSPivoted)
	}
	fuseOverlayFSOption := fmt.Sprintf("lowerdir=%s:%s,upperdir=%s,workdir=%s,squash_to_root",
		tmpRootFS+lowerDir,
		tmpRootFS+tmpRootFS,
		tmpRootFS+upperDir,
		tmpRootFS+workDir)
	fuseOverlayFSArgs := []string{"-o", fuseOverlayFSOption, tmpRootFS + mergedDir}
	if utils.GlobalFlag.FuseOverlayFSArgs != "" {
		fuseOverlayFSArgs = append(fuseOverlayFSArgs, strings.Split(utils.GlobalFlag.FuseOverlayFSArgs, " ")...)
	}
	if utils.GlobalFlag.FuseOverlayFS != "" {
		err = utils.RunCommand(utils.GlobalFlag.FuseOverlayFS, fuseOverlayFSArgs...)
	} else {
		err = utils.ExecFuseOvlMain(fuseOverlayFSArgs)
	}
	if err != nil {
		utils.ExitWith(err, "fuse-overlayfs:", utils.GlobalFlag.FuseOverlayFS, fuseOverlayFSArgs)
	}
	err = utils.MountAll([]utils.MountOption{
		{
			Source: tmpRootFS + "/dev",
			Target: path.Join(tmpRootFS+mergedDir, "dev"),
		},
		{
			Source: tmpRootFS + "/proc",
			Target: path.Join(tmpRootFS+mergedDir, "proc"),
		},
		{
			Source: tmpRootFS + "/home",
			Target: path.Join(tmpRootFS+mergedDir, "home"),
		},
		{
			Source: tmpRootFS + "/project",
			Target: path.Join(tmpRootFS+mergedDir, "project"),
		},
		{
			Source: tmpRootFS + "/root",
			Target: path.Join(tmpRootFS+mergedDir, "root"),
		},
		{
			Source: tmpRootFS + "/tmp",
			Target: path.Join(tmpRootFS+mergedDir, "tmp"),
		},
		{
			Source: tmpRootFS + "/sys",
			Target: path.Join(tmpRootFS+mergedDir, "sys"),
		},
		{
			Source: tmpRootFS + aptDataDir,
			Target: path.Join(tmpRootFS+mergedDir, "/var/lib/apt"),
		},
		{
			Source: tmpRootFS + aptCacheDir,
			Target: path.Join(tmpRootFS+mergedDir, "/var/cache"),
		},
	})
	if err != nil {
		utils.ExitWith(err, "MountAll2")
	}
	err = utils.SwitchTo("MountOverlayStage2", &utils.SwitchFlags{Cloneflags: syscall.CLONE_NEWNS})
	if err != nil {
		utils.ExitWith(err)
	}
}

func BuildMain(cmd *cobra.Command, args []string) error {
	BuildFlag.Args = args
	utils.GlobalFlag.FuseOverlayFS = BuildFlag.FuseOverlayFS
	utils.GlobalFlag.FuseOverlayFSArgs = BuildFlag.FuseOverlayFSArgs
	reexec.Register("MountOverlay", MountOverlay)
	reexec.Register("MountOverlayStage2", MountOverlayStage2)
	if !reexec.Init() {
		if BuildFlag.Strict {
			encodedArgs := []string{}
			args := GetBuildArgs()
			for _, str := range args {
				encoded := base64.StdEncoding.EncodeToString([]byte(str))
				encodedArgs = append(encodedArgs, encoded)
			}
			extArgs := []string{"ll-builder", "run", "--exec", fmt.Sprintf("%s build --encoded-args %s", BuildFlag.Self, strings.Join(encodedArgs, ","))}
			utils.Exec(extArgs...)
		}
		if BuildFlag.EncodedArgs != "" {
			args := []string{}
			for _, item := range strings.Split(BuildFlag.EncodedArgs, ",") {
				value, err := base64.StdEncoding.DecodeString(item)
				if err != nil {
					utils.ExitWith(err)
				}
				args = append(args, string(value))
			}

			args = append([]string{"build"}, args...)
			return utils.SwitchTo("MountOverlay", &utils.SwitchFlags{
				UID:           0,
				GID:           0,
				Cloneflags:    syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
				Args:          args,
				NoDefaultArgs: true,
			})
		} else {
			args := GetBuildArgs()
			args = append([]string{"build"}, args...)
			return utils.SwitchTo("MountOverlay", &utils.SwitchFlags{
				UID:           0,
				GID:           0,
				Cloneflags:    syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
				Args:          args,
				NoDefaultArgs: true,
			})

		}
	}
	return nil
}

func CreateBuildCommand() *cobra.Command {
	cwd, err := os.Getwd()
	if err != nil {
		utils.ExitWith(err)
	}
	execPath, err := os.Executable()
	if err != nil {
		utils.ExitWith(err)
	}

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "进入构建环境",
		Long:    utils.BuildHelpMessage(BuildCommandDescription),
		Example: utils.BuildHelpMessage(BuildCommandHelp),
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(BuildMain(cmd, args))
		},
	}
	cmd.Flags().StringVar(&BuildFlag.RootFS, "rootfs", "/run/host/rootfs", "主机根目录路径")
	cmd.Flags().StringVar(&BuildFlag.TmpRootFS, "tmp-rootfs", "/tmp/rootfs", "临时根目录路径")
	cmd.Flags().StringVar(&BuildFlag.CWD, "cwd", cwd, "当前工作目录路径")
	if _ptrace.IsSupported {
		cmd.Flags().BoolVar(&BuildFlag.Ptrace, "ptrace", false, "修正系统调用(chown)")
	}
	cmd.Flags().StringVar(&BuildFlag.EncodedArgs, "encoded-args", "", "编码后的参数")
	cmd.Flags().StringVar(&BuildFlag.Self, "self", execPath, "github.com/System233/ll-killer-go路径")
	cmd.Flags().BoolVarP(&BuildFlag.Strict, "strict", "x", os.Getenv("LINGLONG_APPID") == "", "严格模式，启动一个与运行时环境相同的构建环境，确保环境一致性（不含gcc等工具）")
	cmd.Flags().StringVar(&BuildFlag.FuseOverlayFS, "fuse-overlayfs", "", "外部fuse-overlayfs命令路径(可选)")
	cmd.Flags().StringVar(&BuildFlag.FuseOverlayFSArgs, "fuse-overlayfs-args", "", "fuse-overlayfs命令额外参数")

	cmd.Flags().MarkHidden("encoded-args")
	// cmd.Flags().MarkHidden("self")
	// cmd.Flags().MarkHidden("cwd")
	return cmd
}
