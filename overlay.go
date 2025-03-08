/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package main

import (
	"os"
	"runtime"
	"unsafe"

	"github.com/spf13/cobra"
)

/*
#cgo pkg-config: fuse3
#cgo LDFLAGS: -L. -lfuse-overlayfs -lgnu
#include <stdlib.h>
#include <stdio.h>
extern int fuse_ovl_main(int argc, char** argv);
int test(int argc, char** argv){
printf("fuse_ovl_main:%p,%d,%s\n",fuse_ovl_main,argc,argv[0]);
return fuse_ovl_main(argc,argv);
}
*/
import "C"

var OverlayFlag struct {
	Self  string
	Shell string
	Args  []string
}

func OverlayMain(cmd *cobra.Command, args []string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	argc := len(args) + 1
	argv := make([]*C.char, argc+1)
	args = append([]string{os.Args[0]}, args...)
	for i, arg := range args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}
	C.fuse_ovl_main(C.int(argc), (**C.char)(unsafe.Pointer(&argv[0])))
	return nil
}

func CreateOverlayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "overlay",
		Short:              "内置fuse-overlayfs挂载",
		Long:               "此命令调用内嵌的fuse-overlayfs程序实现overlay挂载",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			ExitWith(OverlayMain(cmd, args))
		},
	}

	return cmd
}
