/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/System233/ll-killer-go/config"
	"golang.org/x/sys/unix"

	"github.com/moby/sys/reexec"
)

var GlobalFlag struct {
	Debug             bool
	FuseOverlayFS     string
	FuseOverlayFSArgs string
}

type SwitchFlags struct {
	UID           int
	GID           int
	Cloneflags    uintptr
	Args          []string
	NoDefaultArgs bool
	UidMappings   []string
	GidMappings   []string
}
type ExitStatus struct {
	ExitCode int
}

func (s *ExitStatus) Error() string {
	return fmt.Sprint("exited:", s.ExitCode)
}
func RemountProc() error {
	if err := Mount(&MountOption{
		Source: "proc",
		Target: "/proc",
		FSType: "proc",
		Flags:  syscall.MS_RELATIME | syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID,
	}); err != nil {
		return fmt.Errorf("remount:/proc: %v", err)
	}
	return nil
}
func CreateCommand(name string) *exec.Cmd {
	cmd := reexec.Command(name)
	cmd.Args = append(cmd.Args, os.Args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func SwitchTo(next string, flags *SwitchFlags) error {
	cmd := CreateCommand(next)
	if flags.NoDefaultArgs {
		cmd.Args = []string{cmd.Args[0]}
	}
	if len(flags.Args) > 0 {
		cmd.Args = append(cmd.Args, flags.Args...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Unshareflags: flags.Cloneflags,
	}
	if flags.Cloneflags&syscall.CLONE_NEWUSER != 0 {
		if os.Getuid() != flags.UID && os.Getgid() != flags.GID {
			cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{
				{
					ContainerID: flags.UID,
					HostID:      os.Getuid(),
					Size:        1,
				},
			}
			cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{
				{
					ContainerID: flags.GID,
					HostID:      os.Getgid(),
					Size:        1,
				},
			}
		} else {
			cmd.SysProcAttr.Unshareflags ^= syscall.CLONE_NEWUSER
		}
	}
	cmd.SysProcAttr.Cloneflags = flags.Cloneflags ^ syscall.CLONE_NEWNS
	cmd.SysProcAttr.Unshareflags = flags.Cloneflags & syscall.CLONE_NEWNS
	Debug("SwitchTo", fmt.Sprintf("%#x", cmd.SysProcAttr.Unshareflags), cmd.Path, cmd.Args)
	return cmd.Run()
}

func NewCommand(name string, args ...string) *exec.Cmd {
	Debug("NewCommand", name, args)
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
func RunCommand(name string, args ...string) error {
	Debug("RunCommand", name, args)
	cmd := NewCommand(name, args...)
	return cmd.Run()
}

func DefaultShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "/bin/bash"
	}
	return shell
}
func ExecRaw(args ...string) error {
	Debug("ExecRaw", args)
	if len(args) == 0 {
		args = []string{DefaultShell()}
	}

	binary, err := exec.LookPath(args[0])
	if err != nil {
		return fmt.Errorf("%s: %s", binary, err)
	}

	if err := syscall.Exec(binary, args, os.Environ()); err != nil {
		return fmt.Errorf("%s: %s", binary, err)
	}
	return nil
}
func Exec(args ...string) {
	Debug("Exec", args)
	err := ExecRaw(args...)
	if err != nil {
		ExitWith(err)
	}
}
func IsExist(name string) bool {
	_, err := os.Lstat(name)
	return !os.IsNotExist(err)
}

func Mount(opt *MountOption) error {
	Debug("Mount", []string{opt.Source, opt.Target, opt.FSType, opt.FSType, opt.Data})
	if opt.FSType == "" && (opt.Flags == 0 || opt.Flags&syscall.MS_BIND != 0) {
		return MountBind(opt.Source, opt.Target, opt.Flags)
	}
	if opt.FSType == "merge" {
		filesystem := make(map[string]string)
		excludes := append([]string{opt.Target}, strings.Split(opt.Data, config.MountArgsItemSep)...)
		if opt.Data == "" {
			excludes = append(excludes, "/tmp", "/proc", "/dev", "/sys", "/run", "/var/run")
		}
		for index, path := range excludes {
			result, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			excludes[index] = result
		}
		target, err := filepath.Abs(opt.Target)
		if err != nil {
			return err
		}
		sources := strings.Split(opt.Source, config.MountArgsItemSep)
		for index, path := range sources {
			result, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			sources[index] = result
		}
		filesystem, err = opt.buildFileSystem(sources, target, filesystem, excludes)
		if err != nil {
			return err
		}

		for mount, from := range filesystem {
			if err = MountBind(from, mount, opt.Flags); err != nil {
				return err
			}
		}
		return nil
	}
	if opt.FSType == config.FuseOverlayFSType || opt.FSType == "ifovl" {
		fuseOverlayFSArgs := []string{"-o", opt.Data, opt.Target}
		if GlobalFlag.FuseOverlayFSArgs != "" {
			fuseOverlayFSArgs = append(fuseOverlayFSArgs, strings.Split(GlobalFlag.FuseOverlayFSArgs, " ")...)
		}
		dirs := []string{opt.Target}
		for _, item := range strings.Split(opt.Data, ",") {
			param := strings.SplitN(item, "=", 2)
			if len(param) == 2 {
				key := strings.TrimSpace(param[0])
				value := strings.TrimSpace(param[1])
				if key == "upperdir" || key == "workdir" {
					dirs = append(dirs, value)
				}
			}
		}
		Debug("mkdir.overlay", GlobalFlag.FuseOverlayFS, dirs)
		if err := MkdirAlls(dirs, 0755); err != nil {
			log.Println(err)
		}
		if opt.FSType == "ifovl" || GlobalFlag.FuseOverlayFS == "" {
			return ExecFuseOvlMain(fuseOverlayFSArgs)
		}
		return RunCommand(GlobalFlag.FuseOverlayFS, fuseOverlayFSArgs...)
	}
	if opt.FSType != "" {
		err := os.MkdirAll(opt.Target, 0755)
		if err != nil {
			return err
		}
	}
	return syscall.Mount(opt.Source, opt.Target, opt.FSType, uintptr(opt.Flags), opt.Data)
}
func MountBind(source string, target string, flags int) error {
	Debug("MountBind", source, target, flags)
	srcInfo, err := os.Lstat(source)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if srcInfo.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(source)
		if err != nil {
			return err
		}
		dirname := path.Dir(target)
		err = os.MkdirAll(dirname, 0755)
		if err != nil {
			return err
		}
		os.Remove(target)
		err = os.Symlink(link, target)
		if err == nil {
			return nil
		}
		if !os.IsExist(err) {
			return err
		}
	}

	if !IsExist(target) {
		if srcInfo.IsDir() {
			err = os.MkdirAll(target, 0755)
			if err != nil {
				return err
			}
		} else {
			dirname := path.Dir(target)
			err = os.MkdirAll(dirname, 0755)
			if err != nil {
				return err
			}

			_, err = os.Create(target)
			if err != nil {
				return err
			}
		}
	}
	flags |= unix.MS_BIND | unix.MS_PRIVATE
	if srcInfo.IsDir() {
		flags |= unix.MS_REC
	}
	err = unix.Mount(source, target, "none", uintptr(flags), "")
	if err != nil {
		return fmt.Errorf("mount:%s->%s(%#x):%w", source, target, flags, err)
	}
	return nil
}
func MkdirAlls(dirs []string, mode os.FileMode) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

type MountOption struct {
	Source string
	Target string
	FSType string
	Flags  int
	Data   string
}

func MountAll(opts []MountOption) error {
	for _, opt := range opts {
		if err := opt.Mount(); err != nil {
			return fmt.Errorf("mount:%s:%v", opt.Target, err)
		}
	}
	return nil
}

var MountFlagMap = map[string]int{
	"active":      syscall.MS_ACTIVE,
	"async":       syscall.MS_ASYNC,
	"bind":        syscall.MS_BIND,
	"rbind":       syscall.MS_BIND | syscall.MS_REC,
	"dirsync":     syscall.MS_DIRSYNC,
	"invalidate":  syscall.MS_INVALIDATE,
	"i_version":   syscall.MS_I_VERSION,
	"kernmount":   syscall.MS_KERNMOUNT,
	"mandlock":    syscall.MS_MANDLOCK,
	"move":        syscall.MS_MOVE,
	"noatime":     syscall.MS_NOATIME,
	"nodev":       syscall.MS_NODEV,
	"nodiratime":  syscall.MS_NODIRATIME,
	"noexec":      syscall.MS_NOEXEC,
	"nosuid":      syscall.MS_NOSUID,
	"nouser":      syscall.MS_NOUSER,
	"posixacl":    syscall.MS_POSIXACL,
	"private":     syscall.MS_PRIVATE,
	"rdonly":      syscall.MS_RDONLY,
	"rec":         syscall.MS_REC,
	"relatime":    syscall.MS_RELATIME,
	"remount":     syscall.MS_REMOUNT,
	"shared":      syscall.MS_SHARED,
	"silent":      syscall.MS_SILENT,
	"slave":       syscall.MS_SLAVE,
	"strictatime": syscall.MS_STRICTATIME,
	"sync":        syscall.MS_SYNC,
	"synchronous": syscall.MS_SYNCHRONOUS,
	"unbindable":  syscall.MS_UNBINDABLE,
}

func ParseMountFlag(flag string) int {
	flags := strings.Split(flag, config.MountArgsItemSep)
	value := 0
	for _, flag := range flags {
		if flag == "" {
			continue
		}
		item, ok := MountFlagMap[flag]
		if !ok {
			log.Println("ignore mount flag:", flag)
			continue
		}
		value |= item
	}
	return value
}
func ParseMountOption(item string) MountOption {
	chunks := strings.SplitN(item, config.MountArgsSep, 5)
	opt := MountOption{}
	if len(chunks) >= 2 {
		opt.Source = chunks[0]
		opt.Target = chunks[1]
	}
	if len(chunks) >= 3 {
		opt.Flags = ParseMountFlag(chunks[2])
	}
	if len(chunks) >= 4 {
		opt.FSType = chunks[3]
	}
	if len(chunks) >= 5 {
		opt.Data = chunks[4]
	}
	return opt
}

func (opt *MountOption) buildFileSystem(sources []string, target string, filesystem map[string]string, excludes []string) (map[string]string, error) {
	conflicts := make(map[string]struct{})
	for _, source := range sources {
		for _, prefix := range excludes {
			if source == prefix {
				Debug("mount.skip", source, target)
				filesystem[target] = source
				return filesystem, nil
			}
		}
	}

	items := make([][]os.DirEntry, len(sources))
	for i, dir := range sources {
		files, _ := os.ReadDir(dir)
		items[i] = files
	}

	for i, parent := range sources {
		files := items[i]
		if files == nil {
			continue
		}

		for _, file := range files {
			path := filepath.Join(target, file.Name())
			current := filepath.Join(parent, file.Name())
			if _, exists := conflicts[file.Name()]; exists {
				continue
			}
			if _, exists := filesystem[path]; exists && file.IsDir() {
				conflicts[file.Name()] = struct{}{}
				delete(filesystem, path)
				continue
			}
			filesystem[path] = current
		}
	}

	for name := range conflicts {
		var walks []string
		for i, dir := range sources {
			if items[i] != nil {
				next := filepath.Join(dir, name)
				walks = append(walks, next)
			}
		}
		_, err := opt.buildFileSystem(walks, filepath.Join(target, name), filesystem, excludes)
		if err != nil {
			return nil, err
		}
	}
	return filesystem, nil
}

func (opt *MountOption) Mount() error {
	return Mount(opt)
}

func (opt *MountOption) IsFuse() bool {
	return opt.FSType == config.FuseOverlayFSType || opt.FSType == "ifovl"
}
func (opt *MountOption) Unmount() error {
	if opt.FSType != "" && opt.FSType != "merge" {
		err := unix.Unmount(opt.Target, unix.MNT_DETACH)
		Debug("Unmount", opt.Target, err)
		return err
	}
	return nil
}
func Debug(v ...any) {
	if GlobalFlag.Debug {
		log.Println(v...)
	}
}
func ExitWith(err error, v ...any) {
	Debug("ExitWith", err, v)
	if err := WaitForChild(); err != nil {
		Debug("WaitForChild", err)
	}
	if err == nil {
		os.Exit(0)
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	}
	var status *ExitStatus
	if errors.As(err, &status) {
		os.Exit(status.ExitCode)
	}
	log.Fatalln(append([]any{err}, v...)...)
}

func OpenFile(destPath string, perm os.FileMode, force bool) (*os.File, error) {
	flag := 0
	if !force {
		flag = os.O_EXCL
	}
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|flag, 0755)
	if err != nil {
		return nil, err
	}
	return destFile, nil
}
func WriteFile(destPath string, data []byte, perm os.FileMode, force bool) error {
	destFile, err := OpenFile(destPath, 0755, force)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = destFile.Write(data)
	return err
}

func CopyFile(destPath string, src io.Reader, perm os.FileMode, force bool) error {
	dst, err := OpenFile(destPath, 0755, force)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}
func CopySymlink(destPath string, src string, force bool) error {
	if force && IsExist(destPath) {
		os.Remove(destPath)
	}
	err := os.Symlink(src, destPath)
	return err
}
func BuildHelpMessage(help string) string {
	return strings.ReplaceAll(help, "<program>", os.Args[0])
}
func GetKillerExec() (string, error) {
	path, err := os.Executable()
	if err != nil {
		path, err = filepath.Abs(os.Args[0])
		if err != nil {
			return "", err
		}
		os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), filepath.Dir(path)))
	}
	return path, nil
}
func SetupEnvVar() error {
	path, err := GetKillerExec()
	if err != nil {
		return err
	}
	os.Setenv(config.KillerExecEnv, path)
	return nil
}
func CopyFileIO(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func IsSameFile(file1 string, file2 string) (bool, error) {
	stat1, err := os.Lstat(file1)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	stat2, err := os.Lstat(file2)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	sys1 := stat1.Sys().(*syscall.Stat_t)
	sys2 := stat2.Sys().(*syscall.Stat_t)

	return sys1.Ino == sys2.Ino && sys1.Dev == sys2.Dev, nil
}

func Must(err error, v ...any) {
	if err != nil {
		ExitWith(err, v...)
	}
}
