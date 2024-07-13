// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import "fmt"

func panicf(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	s = "babycli: " + s
	panic(s)
}
