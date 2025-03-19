/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package utils

import (
	"context"
	"fmt"
	"log"
	"math"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var childSubreaperWaitDuration time.Duration

func InitChildSubreaper() error {
	_, _, errno := unix.Syscall(syscall.SYS_PRCTL, unix.PR_SET_CHILD_SUBREAPER, 1, 0)
	if errno != 0 {
		return fmt.Errorf("prctl failed: %v", errno)
	}
	_, _, errno = unix.Syscall(syscall.SYS_PRCTL, unix.PR_SET_PDEATHSIG, uintptr(unix.SIGTERM), 0)
	if errno != 0 {
		return fmt.Errorf("prctl failed: %v", errno)
	}
	return nil
}
func SetChildSubreaperWaitDuration(duration time.Duration) {
	if duration < 0 {
		duration = math.MaxInt64
	}
	childSubreaperWaitDuration = duration
}
func WaitForChild() error {
	Debug("WaitForChild")
	duration := childSubreaperWaitDuration
	sig := make(chan error)
	if duration != 0 {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			for {
				select {
				case <-ctx.Done():
					sig <- nil
					return
				default:
					var status syscall.WaitStatus
					pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
					if err != nil {
						if err == unix.ECHILD {
							Debug("所有后台进程已退出")
							sig <- nil
							return
						}
						log.Println("Wait failed: ", err)
						sig <- err
						return
					}
					Debug(fmt.Sprintf("Reaped process %d with status %d\n", pid, status.ExitStatus()))
				}
			}
		}()
		go func() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				log.Println("正在等待所有后台进程退出...")
			}
		}()
	}
	select {
	case err := <-sig:
		return err
	case <-time.After(duration):
		return nil
	}
}
