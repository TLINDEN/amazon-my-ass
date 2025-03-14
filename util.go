/*
Copyright © 2023-2024 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/mattn/go-isatty"
)

func Mkdir(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// returns TRUE if stdout is NOT a tty or windows
func IsNoTty() bool {
	if runtime.GOOS == WIN || !isatty.IsTerminal(os.Stdout.Fd()) {
		return true
	}

	// it is a tty
	return false
}

func GetThrottleTime() time.Duration {
	return time.Duration(rand.Intn(MaxThrottle-MinThrottle+1)+MinThrottle) * time.Millisecond
}

// look if a key in a map exists, generic variant
func Exists[K comparable, V any](m map[K]V, v K) bool {
	if _, ok := m[v]; ok {
		return true
	}

	return false
}
