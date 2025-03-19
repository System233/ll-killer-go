/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package config

const (
	KillerExec        = "ll-killer"
	KillerExecEnv     = "KILLER_EXEC"
	KillerPackerEnv   = "KILLER_PACKER"
	FileSystemDir     = "linglong/filesystem"
	UpperDirName      = "diff"
	LowerDirName      = "overwrite"
	WorkDirName       = "work"
	MergedDirName     = "merged"
	SourceListFile    = "sources.list"
	AptDir            = "linglong/apt"
	AptDataDir        = AptDir + "/data"
	AptCacheDir       = AptDir + "/cache"
	AptDpkgDir        = AptDir + "/dpkg"
	AptConfDir        = "apt.conf.d"
	AptConfFile       = AptConfDir + "/ll-killer.conf"
	LinglongYaml      = "linglong.yaml"
	KillerCommands    = "KILLER_COMMANDS"
	KillerDebug       = "KILLER_DEBUG"
	MountArgsSep      = ":"
	MountArgsItemSep  = "+"
	FuseOverlayFSType = "fuse-overlayfs"
	Repo              = "System233/ll-killer-go"
	GithubURL         = "https://github.com/" + Repo
)

var (
	Version   = "unknown"
	BuildTime = "unknown"
	Tag       = "unknown"
	Variant   = "unknown"
)
