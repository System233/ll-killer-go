/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _apt

import (
	"ll-killer/config"
	"ll-killer/utils"
	"os"
	"syscall"

	"github.com/moby/sys/reexec"
	"github.com/spf13/cobra"
)

const APTCommandHelp = `
此命令用于在宿主机上创建隔离的APT环境。
你可以在隔离环境中使用apt-file，或内置的ldd-search.sh等工具在指定的APT仓库中查找依赖。
当前目录下的apt.conf,apt.conf.d,sources.list和sources.list.d将被挂载至/etc，你可以在这些文件或文件夹中自定义你的apt配置。
隔离环境中的apt缓存将被构建容器重用。

# 示例
<program> apt -- bash
`

var APTFlag struct {
	Args []string
}

func MountAPT() {
	/*
			mkdir -p sources.list.d "$APT_TMP_DIR/apt" "$APT_TMP_DIR/cache"
		    mount --bind ./sources.list /etc/apt/sources.list
		    mount --rbind ./sources.list.d /etc/apt/sources.list.d
		    mount --rbind "$APT_TMP_DIR/apt" /var/lib/apt
		    mount --rbind "$APT_TMP_DIR/cache" /var/cache
		    apt -o APT::Sandbox::User="root" update -y
		    reexec shell
	*/
	err := os.MkdirAll(config.AptDir, 0755)
	if err != nil {
		utils.ExitWith(err)
	}
	err = utils.MkdirAlls([]string{config.AptDataDir, config.AptCacheDir, config.AptConfDir}, 0755)
	if err != nil {
		utils.ExitWith(err)
	}
	err = utils.MountAll([]utils.MountOption{
		{
			Source: "sources.list.d",
			Target: "/etc/apt/sources.list.d",
			Flags:  syscall.MS_BIND,
		},
		{
			Source: "sources.list",
			Target: "/etc/apt/sources.list",
			Flags:  syscall.MS_BIND,
		},
		{
			Source: "auth.conf.d",
			Target: "/etc/apt/auth.conf.d",
			Flags:  syscall.MS_BIND,
		},
		{
			Source: "apt.conf.d",
			Target: "/etc/apt/apt.conf.d",
			Flags:  syscall.MS_BIND,
		},
		{
			Source: "apt.conf",
			Target: "/etc/apt/apt.conf",
			Flags:  syscall.MS_BIND,
		},
		{
			Source: config.AptDataDir,
			Target: "/var/lib/apt",
			Flags:  syscall.MS_BIND,
		},
		{
			Source: config.AptCacheDir,
			Target: "/var/cache",
			Flags:  syscall.MS_BIND,
		},
	})

	if err != nil {
		utils.ExitWith(err, "MountAll")
	}
	utils.Exec(APTFlag.Args...)
}
func APTMain(cmd *cobra.Command, args []string) error {
	APTFlag.Args = args
	reexec.Register("MountAPT", MountAPT)
	if !reexec.Init() {
		return utils.SwitchTo("MountAPT", &utils.SwitchFlags{UID: 0, GID: 0, Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER})
	}
	return nil
}
func CreateAPTCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "apt",
		Short: "进入隔离的APT环境",
		Long:  utils.BuildHelpMessage(APTCommandHelp),
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(APTMain(cmd, args))
		},
	}
}
