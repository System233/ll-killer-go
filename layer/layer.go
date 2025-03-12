package layer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"ll-killer/types"
	"ll-killer/utils"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/sys/unix"
)

const LayerMagic = "<<< deepin linglong layer archive >>>\x00\x00\x00"
const MkfsErofs = "mkfs.erofs"
const DumpErofs = "dump.erofs"
const ErofsFuse = "erofsfuse"
const FuserMount = "fusermount"

type LayerInfo struct {
	Arch          []string `json:"arch"`
	Base          string   `json:"base"`
	Runtime       string   `json:"runtime,omitempty"`
	Channel       string   `json:"channel"`
	Command       []string `json:"command"`
	Description   string   `json:"description"`
	ID            string   `json:"id"`
	Kind          string   `json:"kind"`
	Module        string   `json:"module"`
	Name          string   `json:"name"`
	SchemaVersion string   `json:"schema_version"`
	Size          int64    `json:"size"`
	Version       string   `json:"version"`
}
type LayerInfoHeader struct {
	Info    LayerInfo `json:"info"`
	Version string    `json:"version"`
}

func GetTriplet() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64-linux-gnu"
	case "arm64":
		return "aarch64-linux-gnu"
	case "loong64":
		return "loongarch64-linux-gnu"
	case "mips64":
		return "mips64el-linux-gnuabi64"
	default:
		return "unknown"
	}
}
func getArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "arm64"
	case "loong64":
		return "loongarch64"
	case "sw64":
		return "sw64"
	case "mips64":
		return "mips64"
	default:
		return "unknown"
	}
}
func (info *LayerInfo) ParseLayerInfo(config types.Config) error {
	channel := "main"
	// 验证package版本格式
	if len(strings.Split(config.Package.Version, ".")) != 4 {
		return fmt.Errorf("package.version must be 4-part semver")
	}
	if len(config.Command) == 0 {
		return fmt.Errorf("package.command is empty")
	}

	// 验证kind类型
	if kind := config.Package.Kind; kind != "app" && kind != "runtime" {
		return fmt.Errorf("invalid package.kind: %s", kind)
	}

	// 处理base字段
	normalizedBase, err := normalizeComponent(config.Base, channel, getArch())
	if err != nil {
		return fmt.Errorf("base format error: %w", err)
	}
	info.Base = normalizedBase

	// 处理runtime字段
	if config.Runtime != "" {
		normalizedRuntime, err := normalizeComponent(config.Runtime, channel, getArch())
		if err != nil {
			return fmt.Errorf("runtime format error: %w", err)
		}
		info.Runtime = normalizedRuntime
	}

	// 填充基础字段
	info.Arch = []string{getArch()}
	info.ID = config.Package.ID
	info.Name = config.Package.Name
	info.Version = config.Package.Version
	info.Kind = config.Package.Kind
	info.Description = config.Package.Description
	info.SchemaVersion = "1.0"
	info.Module = "binary"
	info.Channel = channel
	info.Command = config.Command
	return nil
}

func normalizeComponent(input, defaultChannel, defaultArch string) (string, error) {
	// 分割channel部分
	channel := defaultChannel
	if parts := strings.SplitN(input, ":", 2); len(parts) > 1 {
		channel = parts[0]
		input = parts[1]
	}

	// 解析架构和包信息
	parts := strings.Split(input, "/")
	var arch, pkg, version string

	// 检查最后一个元素是否是架构
	if isValidArch(parts[len(parts)-1]) {
		arch = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	} else {
		arch = defaultArch
	}

	// 验证包名和版本
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid component format: %s", input)
	}
	pkg = strings.Join(parts[:len(parts)-1], "/")
	version = parts[len(parts)-1]

	// 验证版本格式
	if !IsValidComponentVersion(version) {
		return "", fmt.Errorf("invalid version format: %s", version)
	}

	return fmt.Sprintf("%s:%s/%s/%s", channel, pkg, version, arch), nil
}

var validArches = map[string]bool{
	"x86_64":      true,
	"arm64":       true,
	"loongarch64": true,
	"sw64":        true,
	"mips64":      true,
}

func isValidArch(arch string) bool {
	return validArches[arch]
}

var versionPartRegex = regexp.MustCompile(`^(0|[1-9]\d*)$`)

func IsValidComponentVersion(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 4 {
		return false
	}
	for _, part := range parts {
		if !versionPartRegex.MatchString(part) {
			return false
		}
	}
	return true
}

func NewLayerInfoHeader(config types.Config) LayerInfoHeader {
	var header LayerInfoHeader
	header.Version = "1"
	header.Info.ParseLayerInfo(config)
	return header

}
func (l *LayerInfoHeader) Print() {
	fmt.Println("Layer元数据版本:")
	fmt.Printf("    版本: %s\n", l.Version)
	l.Info.Print()
}

func (l *LayerInfo) FileName() string {
	return fmt.Sprintf("%s_%s_%s_%s.layer", l.ID, l.Version, strings.Join(l.Arch, "+"), l.Module)
}

func (info *LayerInfo) Print() {
	fmt.Println("Layer元数据:")
	fmt.Printf("    名称: %s\n", info.Name)
	fmt.Printf("    ID: %s\n", info.ID)
	fmt.Printf("    版本: %s\n", info.Version)
	fmt.Printf("    模块: %s\n", info.Module)
	fmt.Printf("    类型: %s\n", info.Kind)
	fmt.Printf("    基础: %s\n", info.Base)
	if info.Runtime != "" {
		fmt.Printf("    运行时: %s\n", info.Runtime)
	}
	fmt.Printf("    渠道: %s\n", info.Channel)
	fmt.Printf("    命令行: %s\n", info.Command)
	fmt.Printf("    元数据版本: %s\n", info.SchemaVersion)
	fmt.Printf("    描述: %s\n", info.Description)
	if info.Size != 0 {
		fmt.Printf("    大小: %d 字节\n", info.Size)
	}
	fmt.Printf("    架构: %s\n", strings.Join(info.Arch, ", "))
}

type LayerHeader struct {
	FileName string
	FileSize int64
	Magic    string
	Info     LayerInfoHeader
	InfoSize int
}

func (l *LayerHeader) DataOffset() int {
	return 40 + 4 + l.InfoSize
}
func (l *LayerHeader) PrintAll() error {
	l.Print()
	l.Info.Print()
	return l.PrintErofs(nil)
}
func (l *LayerHeader) Print() {
	fmt.Println("Layer文件头:")
	fmt.Printf("  文件名: %s\n", l.FileName)
	fmt.Printf("  文件大小: %d 字节\n", l.FileSize)
	fmt.Printf("  魔数: %s\n", l.Magic)
	fmt.Printf("  元数据大小: %d 字节\n", l.InfoSize)
	fmt.Printf("  数据偏移量: %d\n", l.DataOffset())
}

type DumpErofsOption struct {
	Args     []string
	ExecPath string
}

func (l *LayerHeader) PrintErofs(opt *DumpErofsOption) error {
	execPath := DumpErofs
	args := []string{fmt.Sprint("--offset=", l.DataOffset()), l.FileName}
	if opt != nil {
		if opt.ExecPath != "" {
			execPath = opt.ExecPath
		}
		args = append(args, opt.Args...)
	}
	fmt.Println("Erofs信息:")
	return utils.RunCommand(execPath, args...)
}

func NewLayerHeader(file *os.File) (*LayerHeader, error) {
	var layer LayerHeader
	file.Seek(0, io.SeekStart)
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	layer.FileSize = info.Size()
	layer.FileName = info.Name()

	magic := make([]byte, 40)
	_, err = io.ReadFull(file, magic)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(magic, []byte(LayerMagic)) {
		return nil, fmt.Errorf("错误的文件头：%v %v", magic, []byte(LayerMagic))
	}
	layer.Magic = string(magic)

	var metadataSize int32
	err = binary.Read(file, binary.LittleEndian, &metadataSize)
	if err != nil {
		return nil, err
	}
	layer.InfoSize = int(metadataSize)

	metadataBytes := make([]byte, metadataSize)
	_, err = io.ReadFull(file, metadataBytes)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(metadataBytes, &layer.Info)
	if err != nil {
		return nil, err
	}

	return &layer, nil
}
func NewLayerHeaderFromFile(filePath string) (*LayerHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	layer, err := NewLayerHeader(file)
	layer.FileName = filePath
	return layer, err
}

type MountOption struct {
	Source   string
	Target   string
	Args     []string
	ExecPath string
}

func Mount(opt *MountOption) error {
	execPath := ErofsFuse
	if opt.ExecPath != "" {
		execPath = opt.ExecPath
	}
	header, err := NewLayerHeaderFromFile(opt.Source)
	if err != nil {
		return err
	}
	err = os.MkdirAll(opt.Target, 0755)
	if err != nil {
		return err
	}
	args := []string{fmt.Sprint("--offset=", header.DataOffset()), opt.Source, opt.Target}
	if len(opt.Args) > 0 {
		args = append(args, opt.Args...)
	}
	return utils.RunCommand(execPath, args...)
}

type UmountOption struct {
	Target   string
	Args     []string
	ExecPath string
}

func Umount(opt *UmountOption) error {
	execPath := FuserMount
	if opt.ExecPath != "" {
		execPath = opt.ExecPath
	}
	args := []string{"-u", opt.Target}
	if len(opt.Args) > 0 {
		args = append(args, opt.Args...)
	}
	return utils.RunCommand(execPath, args...)
}

type PackOption struct {
	Source     string
	Target     string
	ExecPath   string
	Compressor string
	BlockSize  int
	Uid        int
	Gid        int
	Args       []string
}

func Pack(opt *PackOption) error {
	execPath := MkfsErofs
	var layer LayerHeader
	layer.Info.Version = "1"
	infoJson, err := os.Open(path.Join(opt.Source, "info.json"))
	if err != nil {
		return err
	}
	infoData, err := io.ReadAll(infoJson)
	defer infoJson.Close()
	err = json.Unmarshal(infoData, &layer.Info.Info)
	if err != nil {
		return err
	}
	target := opt.Target
	if target == "" {
		target = layer.Info.Info.FileName()
	}
	if opt.ExecPath != "" {
		execPath = opt.ExecPath
	}
	metadata, err := json.Marshal(layer.Info)
	if err != nil {
		return err
	}
	layer.InfoSize = len(metadata)
	args := []string{}
	args = append(args, opt.Args...)
	if opt.Compressor != "" {
		args = append(args, fmt.Sprint("-z", opt.Compressor))
	}

	sizeData := make([]byte, 4)
	_, err = binary.Encode(sizeData, binary.LittleEndian, int32(layer.InfoSize))
	if err != nil {
		return err
	}
	if opt.Uid >= 0 {
		args = append(args, fmt.Sprint("--force-uid=", opt.Uid))
	}
	if opt.Gid >= 0 {
		args = append(args, fmt.Sprint("--force-gid=", opt.Gid))
	}
	args = append(args,
		// fmt.Sprint("--offset=", layer.DataOffset()), # BUG
		fmt.Sprint("-b", opt.BlockSize),
		target,
		opt.Source)
	err = utils.RunCommand(execPath, args...)
	if err != nil {
		return err
	}

	header := append([]byte(LayerMagic), append(sizeData, metadata...)...)
	fp, err := os.OpenFile(target, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()

	fstat, err := fp.Stat()
	if err != nil {
		return err
	}
	rawSize := fstat.Size()
	newSize := rawSize + int64(layer.DataOffset())
	err = fp.Truncate(newSize)
	if err != nil {
		return err
	}
	mem, err := unix.Mmap(int(fp.Fd()), 0, int(newSize), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return err
	}
	defer unix.Munmap(mem)
	err = unix.Madvise(mem, unix.MADV_SEQUENTIAL)
	if err != nil {
		return err
	}
	if count := copy(mem[layer.DataOffset():], mem); count != int(rawSize) {
		return fmt.Errorf("复制错误:%d!=%d", count, rawSize)
	}

	if count := copy(mem[:layer.DataOffset()], header); count != int(layer.DataOffset()) {
		return fmt.Errorf("复制错误:%d!=%d", count, layer.DataOffset())
	}
	err = unix.Msync(mem, unix.MS_SYNC)
	if err != nil {
		return err
	}
	fmt.Println("文件已输出至:", target)
	return err
}

func Dump(target string) error {
	header, err := NewLayerHeaderFromFile(target)
	if err != nil {
		return err
	}
	return header.PrintAll()
}

func NormalizeVersion(version string) string {
	re := regexp.MustCompile(`\D`)
	chunks := strings.SplitN(version, ".", 4)
	for index, chunk := range chunks {
		version := strings.TrimLeft(strings.TrimSpace(re.ReplaceAllString(chunk, "")), "0")
		if version == "" {
			version = "0"
		}
		chunks[index] = version
	}
	for len(chunks) < 4 {
		chunks = append(chunks, "0")
	}
	return strings.Join(chunks, ".")
}
