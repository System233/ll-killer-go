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
	"text/template"

	"github.com/System233/ll-killer-go/apps/layer/fs"
	"github.com/System233/ll-killer-go/config"
	"github.com/System233/ll-killer-go/layer"
	"github.com/System233/ll-killer-go/reexec"
	"github.com/System233/ll-killer-go/utils"

	"github.com/spf13/cobra"
)

var Flag struct {
	RootFs      string
	Runtime     string
	Target      string
	ExecPath    string
	Compressor  string
	BlockSize   int
	Gid         int
	Uid         int
	NoPostSetup bool
	PrintJson   bool
	Format      string
	Quiet       bool
	Print       string
	NoLayer     bool
	PackArgs    []string
	Args        []string
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
const WorkDir = "linglong/output"

var Config layer.Config
var LayerInfo layer.LayerInfo

func PostPackUp(workDir string) error {
	baseDir := path.Join(workDir, "layer")
	appID := Config.Package.ID

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}

	// 创建files目录并挂载
	filesDir := path.Join(baseDir, "files")
	buildHostDir := path.Join(workDir, "build")
	if err := os.MkdirAll(filesDir, 0755); err != nil {
		return err
	}
	if err := utils.MountBind(buildHostDir, filesDir, syscall.MS_BIND); err != nil {
		return err
	}

	// 创建entries目录
	entriesDir := path.Join(baseDir, "entries")

	if err := os.MkdirAll(entriesDir, 0755); err != nil {
		return err
	}
	// 处理share目录硬链接
	shareSrc := path.Join(filesDir, "share")
	if _, err := os.Stat(shareSrc); err == nil {
		shareDst := path.Join(entriesDir, "share")
		if err := utils.MountBind(shareSrc, shareDst, 0); err != nil {
			return err
		}
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

	if err := os.WriteFile(installFile, []byte(fileList.String()), 0644); err != nil {
		return err
	}

	LayerInfo.Size = totalSize
	// 复制linglong.yaml
	yamlSrc := path.Join(".", "linglong.yaml")
	yamlDst := path.Join(baseDir, "linglong.yaml")

	if err := utils.CopyFileIO(yamlSrc, yamlDst); err != nil {
		return err
	}

	infoJsonPath := path.Join(baseDir, "info.json")
	infoJsonData, err := json.Marshal(LayerInfo)

	if err != nil {
		return fmt.Errorf("无法序列化info.json:%v", err)
	}

	if err := os.WriteFile(infoJsonPath, infoJsonData, 0755); err != nil {
		return fmt.Errorf("生成info.json失败:%v", err)
	}
	if err := layer.Pack(&layer.PackOption{
		Source:     baseDir,
		Target:     Flag.Target,
		Compressor: Flag.Compressor,
		BlockSize:  Flag.BlockSize,
		Gid:        Flag.Gid,
		Uid:        Flag.Uid,
		Args:       Flag.PackArgs,
	}); err != nil {
		return fmt.Errorf("生成layer失败:%v", err)
	}
	return nil
}

func RunBuildScript(workDir string) error {
	if len(Flag.Args) == 0 {
		Flag.Args = append(Flag.Args, "linglong/entry.sh")
	}
	cmd := utils.NewCommand(Flag.Args[0], Flag.Args[1:]...)
	cmd.Dir = "/project"

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("构建失败:%v", err)
	}
	return nil
}

func RunPostSetup(workDir string) error {
	cmd := utils.NewCommand(PostSetupScript)
	cmd.Dir = "/project"
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("后处理失败:%v", err)
	}
	return nil
}

func BuildLayer() error {
	log.Println("[准备构建环境]")

	if err := utils.RemountProc(); err != nil {
		return err
	}

	if err := fs.SetupFilesystem(fs.SetupFilesystemOption{
		RootFs:    Flag.RootFs,
		Runtime:   Flag.Runtime,
		Quiet:     Flag.Quiet,
		Config:    &Config,
		LayerInfo: &LayerInfo,
	}); err != nil {
		return err
	}

	log.Println("[运行构建脚本]")
	if err := RunBuildScript(WorkDir); err != nil {
		return err
	}

	if !Flag.NoPostSetup {
		log.Println("[文件后处理]")
		if err := RunPostSetup(WorkDir); err != nil {
			return err
		}
	}
	if err := fs.PostFilesystem(); err != nil {
		return err
	}
	if Flag.NoLayer {
		return nil
	}
	log.Println("[打包输出]")
	return PostPackUp(WorkDir)
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
func PrintLayerInfo(key string, json bool, format string) error {
	var cfg layer.Config
	// var info layer.LayerInfo
	type DumpInfo struct {
		layer.LayerInfo
		Layer string `json:"layer"`
	}
	var info DumpInfo

	if err := utils.LoadYamlFile(config.LinglongYaml, &cfg); err != nil {
		return fmt.Errorf("读取linglong.yaml失败:%v", err)
	}

	if err := info.ParseLayerInfo(cfg); err != nil {
		return fmt.Errorf("linglong.yaml配置不合法:%v", err)
	}
	info.Layer = info.LayerInfo.FileName()
	if Flag.Target != "" {
		info.Layer = Flag.Target
	}
	var mapData map[string]interface{}
	data, err := utils.DumpJsonDataIndent(info, "  ")
	if err != nil {
		return err
	}
	if json {
		fmt.Println(string(data))
		return nil
	}
	if err := utils.LoadJsonData(data, &mapData); err != nil {
		return err
	}
	if format != "" {
		t, err := template.New("format").Parse(Flag.Format)
		if err != nil {
			return fmt.Errorf("解析format失败:%v", err)
		}
		return t.Execute(os.Stdout, mapData)
	}
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
	ok, err := reexec.Init()
	if ok {
		return err
	}
	if Flag.Print != "" || Flag.PrintJson || Flag.Format != "" {
		target := "layer"
		if Flag.Print != "" {
			target = Flag.Print
		}
		return PrintLayerInfo(target, Flag.PrintJson, Flag.Format)
	}
	return utils.SwitchTo("BuildLayer", &utils.SwitchFlags{
		UID:           0,
		GID:           0,
		Cloneflags:    syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID,
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
	cmd.Flags().StringVar(&Flag.Runtime, "runtime", "", "runtime文件系统")
	cmd.Flags().IntVarP(&Flag.BlockSize, "block-size", "b", 4096, "块大小")
	cmd.Flags().StringVar(&Flag.ExecPath, "exec", layer.MkfsErofs, "指定mkfs.erofs命令位置")
	cmd.Flags().StringVarP(&Flag.Compressor, "compressor", "z", "lz4hc", "压缩算法，请查看mkfs.erofs帮助")
	cmd.Flags().IntVarP(&Flag.Uid, "force-uid", "U", os.Getuid(), "文件Uid,-1为不更改")
	cmd.Flags().IntVarP(&Flag.Gid, "force-gid", "G", os.Getegid(), "文件Gid,-1为不更改")
	cmd.Flags().BoolVar(&Flag.NoPostSetup, "no-post-setup", false, "不对构建结果进行后处理")
	cmd.Flags().BoolVar(&Flag.NoLayer, "no-layer", false, "不输出layer文件")
	cmd.Flags().StringVar(&Flag.Print, "print", "", "输出应用参数后退出")
	cmd.Flags().BoolVar(&Flag.PrintJson, "json", false, "输出应用参数的JSON格式后退出")
	cmd.Flags().BoolVar(&Flag.Quiet, "quiet", false, "安静模式，构建前不输出项目信息")
	cmd.Flags().StringVarP(&Flag.Target, "output", "o", "", "输出的layer文件名")
	cmd.Flags().StringVarP(&Flag.Format, "format", "f", "", "格式化输出，字段详见json")
	cmd.Flags().StringSliceVar(&Flag.PackArgs, "erofs-args", []string{}, "其他mkfs.erofs选项,逗号分隔")
	cmd.Flags().SortFlags = false
	return cmd
}
