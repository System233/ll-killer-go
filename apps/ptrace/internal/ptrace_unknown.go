//go:build linux && !amd64 && !arm64 && !loong64

/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */

package internal

import (
	"golang.org/x/sys/unix"
)

const IsSupported = true

func PtraceHandle(pid int, regs unix.PtraceRegs) error {
	return nil
}
