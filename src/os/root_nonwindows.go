// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows

package os

func rootCleanPath(s string, prefix, suffix []string) (string, error) {
	return s, nil
}
