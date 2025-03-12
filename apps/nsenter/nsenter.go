/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package _nsenter

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/System233/ll-killer-go/utils"

	"github.com/moby/sys/reexec"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

/*
#cgo CFLAGS: -Wall
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/ioctl.h>
#include <errno.h>
#include <dirent.h>
#include <linux/unistd.h>
#include <sys/syscall.h>

#define MAX_NS 100
#define DELIM ","
int NsEnter(const char *nsPath,const char *nsList[],int total) {
 	const char*nsenterQuietEnv=getenv("KILLER_NSENTER_QUIET");
    // const char *nsList[] = {"cgroup", "ipc", "uts", "net", "pid", "mnt", "time", "user"};
    int fps[MAX_NS];
    int count = 0;
    int result[MAX_NS] = {0};  // 0: not set, 1: set
    int errors[MAX_NS];

	static char nsFile[512];
    // Open namespace files
    for (int i = 0; i < total; i++) {
        snprintf(nsFile, sizeof(nsFile), "%s/%s", nsPath, nsList[i]);
        fps[i] = open(nsFile, O_RDONLY);
        if (fps[i] == -1) {
            perror("Failed to open namespace");
            return errno;
        }
    }

    // Set namespaces
    for (int i = 0; i < total; i++) {
        for (int j = 0; j < total; j++) {
            if (result[j]) continue;

            // Use syscall to set namespace
            if (syscall(SYS_setns, fps[j], 0) == -1) {
                errors[j] = errno;
            } else {
                count++;
                result[j] = 1;
            }
        }

        if (count == total) {
            break;
        }
    }

    // Check for failures
	if (count != total) {
		if(!nsenterQuietEnv){
				for (int i = 0; i < total; i++) {
					if (!result[i]) {
						fprintf(stderr, "setns '%s' failed: %s\n", nsList[i], strerror(errors[i]));
					}
				}
		}
		return -1;
	}

    // Clean up
    for (int i = 0; i < total; i++) {
        close(fps[i]);
    }

    return 0;
}


void __attribute__((constructor)) init(void) {
 	const char*nsenterEnv=getenv("KILLER_NSENTER");
 	const char*nsenterNsEnv=getenv("KILLER_NSENTER_NS");
 	const char*nsenterNoFailEnv=getenv("KILLER_NSENTER_NF");
	if(nsenterEnv){
		static const char *nsList[MAX_NS];
		static char buffer[1000];
		strncpy(buffer, nsenterNsEnv, sizeof(buffer));
		buffer[sizeof(buffer) - 1] = '\0';

		int count = 0;
		char *token = strtok(buffer, DELIM);
		while (token && count < MAX_NS) {
			nsList[count++] = token;
			token = strtok(NULL, DELIM);
		}
		int ret=NsEnter(nsenterEnv,nsList,count);
		if(nsenterNoFailEnv&&ret!=0){
			exit(1);
		}
	}
}
*/
import "C"

var NsEnterFlag struct {
	Pid      int
	Gid      int
	Uid      int
	NsType   []string
	Keep     bool
	NoFail   bool
	Quiet    bool
	Args     []string
	ProcPath string
}

func GetExecArgs() []string {
	args := []string{
		fmt.Sprint(NsEnterFlag.Pid),
		"--uid", fmt.Sprint(NsEnterFlag.Uid),
		"--gid", fmt.Sprint(NsEnterFlag.Gid),
	}

	if NsEnterFlag.ProcPath != "" {
		args = append(args, "--proc", NsEnterFlag.ProcPath)
	}
	if NsEnterFlag.Keep {
		args = append(args, "--keep")
	}
	if !NsEnterFlag.NoFail {
		args = append(args, "--no-fail=false")
	}
	if NsEnterFlag.Quiet {
		args = append(args, "--quiet=false")
	}
	if len(NsEnterFlag.NsType) > 0 {
		args = append(args, "--nstype", strings.Join(NsEnterFlag.NsType, ","))
	}
	if len(NsEnterFlag.Args) > 0 {
		args = append(args, "--")
		args = append(args, NsEnterFlag.Args...)
	}

	return args
}

func ReadMapping(mapPath string) ([]syscall.SysProcIDMap, error) {
	fp, err := os.Open(mapPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	result := []syscall.SysProcIDMap{}
	for {
		var innerId, outerId, count int
		_, err = fmt.Fscanf(fp, "%d%d%d", &innerId, &outerId, &count)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		result = append(result, syscall.SysProcIDMap{
			ContainerID: innerId,
			HostID:      outerId,
			Size:        count,
		})
	}
	return result, nil
}
func GetUid() (int, error) {
	mapping, err := ReadMapping("/proc/self/uid_map")
	if err != nil {
		return 0, err
	}
	if len(mapping) == 0 {
		return 0, fmt.Errorf("没有有效UID")
	}
	return mapping[0].ContainerID, nil
}
func GetGid() (int, error) {
	mapping, err := ReadMapping("/proc/self/gid_map")
	if err != nil {
		return 0, err
	}
	if len(mapping) == 0 {
		return 0, fmt.Errorf("没有有效GID")
	}
	return mapping[0].ContainerID, nil
}
func NsSetupUser() error {

	if NsEnterFlag.Gid >= 0 {
		err := unix.Setgid(NsEnterFlag.Gid)
		if err != nil {
			return fmt.Errorf("setgid: %v", err)
		}
	} else {
		gid, err := GetGid()
		if err == nil {
			err = unix.Setgid(gid)
		}
		if err != nil {
			fmt.Println("setgid:", err)
		}
	}
	if NsEnterFlag.Uid >= 0 {
		err := unix.Setuid(NsEnterFlag.Uid)
		if err != nil {
			return fmt.Errorf("setuid: %v", err)
		}
	} else {
		gid, err := GetUid()
		if err == nil {
			err = unix.Setuid(gid)
		}
		if err != nil {
			fmt.Println("setuid:", err)
		}
	}
	return nil
}
func NsEnterMain(cmd *cobra.Command, args []string) error {
	nsenterEnv := os.Getenv("KILLER_NSENTER")
	if nsenterEnv != "" {
		os.Setenv("KILLER_NSENTER", "")
		os.Setenv("KILLER_NSENTER_NS", "")
		os.Setenv("KILLER_NSENTER_QUIET", "")
		os.Setenv("KILLER_NSENTER_NF", "")
		if !NsEnterFlag.Keep {
			err := NsSetupUser()
			if err != nil {
				return err
			}
		}
		return utils.ExecRaw(args...)
	}
	os.Setenv("KILLER_NSENTER", path.Join("/proc", fmt.Sprint(NsEnterFlag.Pid), "ns"))
	os.Setenv("KILLER_NSENTER_NS", strings.Join(NsEnterFlag.NsType, ","))
	if NsEnterFlag.NoFail {
		os.Setenv("KILLER_NSENTER_NF", "1")
	}
	if NsEnterFlag.Quiet {
		os.Setenv("KILLER_NSENTER_QUIET", "1")
	}
	return utils.ExecRaw(append([]string{reexec.Self(), "nsenter"}, GetExecArgs()...)...)
}
func NsEnterNsEnterCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "nsenter <pid> -- cmd",
		Short: "进入命名空间",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pid, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				utils.ExitWith(err)
			}
			NsEnterFlag.Pid = int(pid)
			NsEnterFlag.Args = args[1:]
			utils.ExitWith(NsEnterMain(cmd, args[1:]))
		},
	}
	cmd.Flags().StringVar(&NsEnterFlag.ProcPath, "proc", "/proc", "指定/proc路径")
	cmd.Flags().IntVarP(&NsEnterFlag.Gid, "gid", "G", os.Getgid(), "指定命名空间内的gid，-1为自动")
	cmd.Flags().IntVarP(&NsEnterFlag.Uid, "uid", "U", os.Getuid(), "指定命名空间内的uid，-1为自动")
	cmd.Flags().BoolVarP(&NsEnterFlag.Keep, "keep", "k", false, "不要切换UID和GID")
	cmd.Flags().BoolVarP(&NsEnterFlag.Quiet, "quiet", "q", false, "不要输出命名空间错误信息")
	cmd.Flags().BoolVarP(&NsEnterFlag.NoFail, "no-fail", "f", true, "切换命名空间失败时退出")
	cmd.Flags().StringSliceVarP(&NsEnterFlag.NsType, "nstype", "N", []string{"cgroup", "ipc", "uts", "net", "pid", "mnt", "time", "user"}, "切换的命名空间类型")
	cmd.Flags().SortFlags = false
	return cmd
}
