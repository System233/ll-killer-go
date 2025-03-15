/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _build

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/System233/ll-killer-go/config"
	"github.com/System233/ll-killer-go/layer"
	"github.com/System233/ll-killer-go/utils"

	"github.com/moby/sys/reexec"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var Flag struct {
	RootFs         string
	Target         string
	ExecPath       string
	Compressor     string
	BlockSize      int
	Gid            int
	Uid            int
	NoPostSetup    bool
	PrintLayerName bool
	PrintJson      bool
	Print          string
	NoLayer        bool
	PackArgs       []string
	Args           []string
}

const BuildCommandDescription = `无需ll-builder, 直接将当前项目构建为layer。
 此过程在宿主机上进行，提供与ll-builder类似的环境，绕过ll-builder，避免不必要的ostree提交和磁盘复制。
 
 此构建模式提供以下内容：
 
 ## 环境变量
 LINGLONG_APPID="{APPID}"
 PREFIX="/opt/apps/{APPID}/files"
 TRIPLET="x86_64-linux-gnu|aarch64-linux-gnu|loongarch64-linux-gnu|..." 
 KILLER_PACKER=1

 ## 目录
 /project: 项目目录
 /: 与宿主机相同

 ## 后处理
 * 为快捷方式和服务单元添加ll-cli run前缀

 * KILLER_PACKER 标识当前处于killer环境，killer环境下setup.sh会自动跳过符号链接修复，
   可以在启动前设置KILLER_PACKER=0禁用该行为。
 `
const BuildCommandHelp = ``
const PostSetupScript = "build-aux/post-setup.sh"

var Config layer.Config
var LayerInfo layer.LayerInfo

func PostPackUp(workDir string) {
	baseDir := path.Join(workDir, "layer")
	appID := Config.Package.ID

	utils.Must(os.MkdirAll(baseDir, 0755))

	// 创建files目录并挂载
	filesDir := path.Join(baseDir, "files")
	buildHostDir := path.Join(workDir, "build")
	utils.Must(os.MkdirAll(filesDir, 0755))
	utils.Must(utils.MountBind(buildHostDir, filesDir, syscall.MS_BIND))

	// 创建entries目录
	entriesDir := path.Join(baseDir, "entries")
	utils.Must(os.MkdirAll(entriesDir, 0755))

	// 处理share目录硬链接
	shareSrc := path.Join(filesDir, "share")
	if _, err := os.Stat(shareSrc); err == nil {
		shareDst := path.Join(entriesDir, "share")
		utils.Must(utils.MountBind(shareSrc, shareDst, 0))
	}

	// 生成install文件
	installFile := path.Join(baseDir, appID+".install")
	var totalSize int64
	var fileList strings.Builder

	filepath.Walk(filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == filesDir {
			return nil
		}
		relPath := strings.TrimPrefix(path, filesDir)
		fileList.WriteString(relPath)
		fileList.WriteString("\n")
		totalSize += info.Size()
		return nil
	})

	utils.Must(os.WriteFile(installFile, []byte(fileList.String()), 0644))

	LayerInfo.Size = totalSize
	// 复制linglong.yaml
	yamlSrc := path.Join(".", "linglong.yaml")
	yamlDst := path.Join(baseDir, "linglong.yaml")
	utils.Must(utils.CopyFileIO(yamlSrc, yamlDst))

	infoJsonPath := path.Join(baseDir, "info.json")
	infoJsonData, err := json.Marshal(LayerInfo)
	utils.Must(err, "无法序列化info.json")

	utils.Must(os.WriteFile(infoJsonPath, infoJsonData, 0755), "生成info.json失败")

	utils.Must(layer.Pack(&layer.PackOption{
		Source:     baseDir,
		Target:     Flag.Target,
		Compressor: Flag.Compressor,
		BlockSize:  Flag.BlockSize,
		Gid:        Flag.Gid,
		Uid:        Flag.Uid,
		Args:       Flag.PackArgs,
	}), "生成layer失败")
}

func SetupFilesystem(workDir string) {
	err := utils.LoadYamlFile(config.LinglongYaml, &Config)
	if err != nil {
		utils.ExitWith(err)
	}

	utils.Must(LayerInfo.ParseLayerInfo(Config), "解析yaml错误")
	LayerInfo.Print()

	if utils.IsExist(workDir) {
		utils.Must(os.RemoveAll(workDir), "无法移除"+workDir)
	}
	if err := os.MkdirAll(workDir, 0755); err != nil {
		utils.ExitWith(err)
	}

	configID := Config.Package.ID
	configBuild := Config.Build
	rootfsPath := path.Join(workDir, "rootfs")

	// 挂载宿主机根目录到rootfsPath（只读）
	if err := utils.Mount(&utils.MountOption{Source: "/", Target: rootfsPath, FSType: "merge", Flags: unix.MS_RDONLY}); err != nil {
		utils.ExitWith(err, "挂载宿主机根目录失败")
	}

	// 创建并挂载/run/host/rootfs
	runTmpfs := path.Join(rootfsPath, "run")
	if err := utils.Mount(&utils.MountOption{Source: "tmpfs", Target: runTmpfs, FSType: "tmpfs"}); err != nil {
		utils.ExitWith(err, "挂载run目录失败")
	}
	// 创建并挂载/run/host/rootfs
	runHostRootfs := path.Join(rootfsPath, "run/host/rootfs")
	if err := os.MkdirAll(runHostRootfs, 0755); err != nil {
		utils.ExitWith(err, "创建run/host/rootfs目录失败")
	}

	// 挂载当前目录到/project
	projectDir := path.Join(rootfsPath, "project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		utils.ExitWith(err, "创建/project目录失败")
	}

	if err := utils.MountBind(".", projectDir, syscall.MS_BIND); err != nil {
		utils.ExitWith(err, "挂载项目目录失败")
	}

	// 创建并挂载输出目录
	buildHostDir := path.Join(workDir, "build")
	if err := os.MkdirAll(buildHostDir, 0755); err != nil {
		utils.ExitWith(err, "创建输出目录失败")
	}

	// 创建并挂载/run/host/rootfs
	optTmpfs := path.Join(rootfsPath, "opt")
	if err := utils.Mount(&utils.MountOption{Source: "tmpfs", Target: optTmpfs, FSType: "tmpfs"}); err != nil {
		utils.ExitWith(err, "挂载opt目录失败")
	}
	optAppsDir := path.Join(rootfsPath, "opt/apps", configID, "files")
	if err := os.MkdirAll(optAppsDir, 0755); err != nil {
		utils.ExitWith(err, "创建/opt/apps目录失败")
	}
	if err := utils.MountBind(buildHostDir, optAppsDir, syscall.MS_BIND); err != nil {
		utils.ExitWith(err, "挂载输出目录失败")
	}

	// 写入entry.sh
	entryPath := "linglong/entry.sh"
	entryData := []byte(fmt.Sprintf("#!/bin/bash\n%s", configBuild))
	if err := os.WriteFile(entryPath, entryData, 0755); err != nil {
		utils.ExitWith(err, "写入entry.sh失败")
	}

	// PivotRoot
	if err := utils.MountBind(rootfsPath, rootfsPath, 0); err != nil {
		utils.ExitWith(err, "绑定根目录失败")
	}

	if err := syscall.PivotRoot(rootfsPath, runHostRootfs); err != nil {
		utils.ExitWith(err, "切换根目录失败")
	}

	// 配置环境变量
	os.Setenv("LINGLONG_APPID", configID)
	os.Setenv("PREFIX", path.Join("/opt/apps", configID, "files"))
	os.Setenv("TRIPLET", layer.GetTriplet())

}
func RunBuildScript(workDir string) {
	if len(Flag.Args) == 0 {
		Flag.Args = append(Flag.Args, "linglong/entry.sh")
	}
	cmd := utils.NewCommand(Flag.Args[0], Flag.Args[1:]...)
	cmd.Dir = "/project"

	if err := cmd.Run(); err != nil {
		utils.ExitWith(err, "构建失败")
	}
}

func RunPostSetup(workDir string) {
	cmd := utils.NewCommand(PostSetupScript)
	cmd.Dir = "/project"
	if err := cmd.Run(); err != nil {
		utils.ExitWith(err, "后处理失败")
	}
}
func BuildLayer() {
	killerPackerEnv := os.Getenv(config.KillerPackerEnv)
	if killerPackerEnv == "" {
		os.Setenv(config.KillerPackerEnv, "1")
	}
	workDir := "linglong/output"
	log.Println("[准备构建环境]")
	SetupFilesystem(workDir)

	log.Println("[运行构建脚本]")
	RunBuildScript(workDir)

	if !Flag.NoPostSetup {
		log.Println("[文件后处理]")
		RunPostSetup(workDir)
	}

	if !Flag.NoLayer {
		log.Println("[打包输出]")
		PostPackUp(workDir)
	}

}
func GetBuildArgs() []string {
	args := []string{
		fmt.Sprint("--block-size=", Flag.BlockSize),
		fmt.Sprint("--force-gid=", Flag.Gid),
		fmt.Sprint("--force-uid=", Flag.Uid),
		fmt.Sprint("--no-layer=", Flag.NoLayer),
		fmt.Sprint("--no-post-setup=", Flag.NoPostSetup),
	}
	if Flag.Target != "" {
		args = append(args, "--output", Flag.Target)
	}
	if Flag.Compressor != "" {
		args = append(args, "--compressor", Flag.Compressor)
	}
	if Flag.ExecPath != "" {
		args = append(args, "--exec", Flag.ExecPath)
	}
	if Flag.RootFs != "" {
		args = append(args, "--rootfs", Flag.RootFs)
	}
	if len(Flag.PackArgs) > 0 {
		args = append(args, "--erofs-args", strings.Join(Flag.PackArgs, ","))
	}

	if len(Flag.Args) > 0 {
		args = append(args, "--")
		args = append(args, Flag.Args...)
	}
	return args
}
func PrintLayerInfo(key string, json bool) error {
	var cfg layer.Config
	// var info layer.LayerInfo
	type DumpInfo struct {
		layer.LayerInfo
		FileName string `json:"fileName"`
	}
	var info DumpInfo
	if Flag.Target != "" {
		fmt.Println(Flag.Target)
		return nil
	}
	utils.Must(utils.LoadYamlFile(config.LinglongYaml, &cfg), "读取linglong.yaml失败")
	utils.Must(info.ParseLayerInfo(cfg), "linglong.yaml配置不合法")
	info.FileName = info.LayerInfo.FileName()
	var mapData map[string]interface{}
	data, err := utils.DumpJsonData(info)
	utils.Must(err)
	if json {
		fmt.Println(string(data))
		return nil
	}
	utils.Must(utils.LoadJsonData(data, &mapData))
	value, ok := mapData[key]
	if !ok {
		return nil
	}
	fmt.Println(value)
	return nil
}
func BuildMain(cmd *cobra.Command, args []string) error {
	Flag.Args = args
	reexec.Register("BuildLayer", BuildLayer)
	if reexec.Init() {
		return nil
	}
	if Flag.PrintLayerName || Flag.Print != "" || Flag.PrintJson {
		target := "fileName"
		if Flag.Print != "" {
			target = Flag.Print
		}
		return PrintLayerInfo(target, Flag.PrintJson)
	}
	return utils.SwitchTo("BuildLayer", &utils.SwitchFlags{
		UID:           0,
		GID:           0,
		Cloneflags:    syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Args:          append([]string{"layer", "build"}, GetBuildArgs()...),
		NoDefaultArgs: true,
	})
}

func CreateBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build [flags] -- cmd",
		Short:   "无需ll-builder, 直接将当前项目构建为layer。",
		Long:    utils.BuildHelpMessage(BuildCommandDescription),
		Example: utils.BuildHelpMessage(BuildCommandHelp),
		Run: func(cmd *cobra.Command, args []string) {
			utils.ExitWith(BuildMain(cmd, args))
		},
	}
	cmd.Flags().StringVar(&Flag.RootFs, "rootfs", "/", "根文件系统")
	cmd.Flags().IntVarP(&Flag.BlockSize, "block-size", "b", 4096, "块大小")
	cmd.Flags().StringVar(&Flag.ExecPath, "exec", layer.MkfsErofs, "指定mkfs.erofs命令位置")
	cmd.Flags().StringVarP(&Flag.Compressor, "compressor", "z", "lz4hc", "压缩算法，请查看mkfs.erofs帮助")
	cmd.Flags().IntVarP(&Flag.Uid, "force-uid", "U", os.Getuid(), "文件Uid,-1为不更改")
	cmd.Flags().IntVarP(&Flag.Gid, "force-gid", "G", os.Getegid(), "文件Gid,-1为不更改")
	cmd.Flags().BoolVar(&Flag.NoPostSetup, "no-post-setup", false, "不对构建结果进行后处理")
	cmd.Flags().BoolVar(&Flag.NoLayer, "no-layer", false, "不输出layer文件")
	cmd.Flags().BoolVar(&Flag.PrintLayerName, "print-layer-name", false, "输出构建的layer文件名")
	cmd.Flags().StringVar(&Flag.Print, "print", "", "输出应用参数后退出")
	cmd.Flags().BoolVar(&Flag.PrintJson, "json", false, "输出应用参数的JSON格式后退出")
	cmd.Flags().StringVarP(&Flag.Target, "output", "o", "", "输出的layer文件名")
	cmd.Flags().StringSliceVar(&Flag.PackArgs, "erofs-args", []string{}, "其他mkfs.erofs选项,逗号分隔")
	cmd.Flags().SortFlags = false
	return cmd
}
