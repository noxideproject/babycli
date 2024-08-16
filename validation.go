// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"io"
)

func (c *Component) validate(output io.Writer) bool {
	ok := true

	for _, f := range c.Flags {
		if len(f.Long) == 1 {
			writef(output, "babycli: long flag %q must be more than one character", f.Long)
			ok = false
		}
		if len(f.Short) > 1 {
			writef(output, "babycli: short flag %q must be one character", f.Short)
			ok = false
		}
	}

	for _, cmd := range c.Components {
		switch len(cmd.Name) {
		case 0:
			writef(output, "babycli: component name missing")
			ok = false
		case 1:
			writef(output, "babycli: component %q must be more than one character", cmd.Name)
			ok = false
		}
	}

	return ok
}
