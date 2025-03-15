/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package reexec

import (
	"fmt"
	"os"
)

var registeredInitializers = make(map[string]func() error)

func Register(name string, initializer func() error) {
	if _, exists := registeredInitializers[name]; exists {
		panic(fmt.Sprintf("reexec func already registered under name %q", name))
	}

	registeredInitializers[name] = initializer
}

func Init() (bool, error) {
	if initializer, ok := registeredInitializers[os.Args[0]]; ok {
		return true, initializer()
	}
	return false, nil
}
