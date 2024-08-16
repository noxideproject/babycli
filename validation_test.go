// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

func TestComponent_validate_short_flag(t *testing.T) {
	t.Parallel()

	config := &Configuration{
		Top: &Component{
			Flags: Flags{
				{
					Long:  "long",
					Short: "xyz",
				},
			},
		},
	}

	w := new(bytes.Buffer)
	c := New(config)
	c.output = w

	result := c.Run()
	must.One(t, result)
	message := strings.TrimSpace(w.String())
	must.Eq(t, `babycli: short flag "xyz" must be one character`, message)
}

func TestComponent_validate_long_flag(t *testing.T) {
	t.Parallel()

	config := &Configuration{
		Top: &Component{
			Flags: Flags{
				{
					Long:  "x",
					Short: "z",
				},
			},
		},
	}

	w := new(bytes.Buffer)
	c := New(config)
	c.output = w

	result := c.Run()
	must.One(t, result)
	message := strings.TrimSpace(w.String())
	must.Eq(t, `babycli: long flag "x" must be more than one character`, message)
}

func TestComponent_validate_name_empty(t *testing.T) {
	t.Parallel()

	config := &Configuration{
		Top: &Component{
			Components: Components{
				{
					Name: "first",
				},
				{
					Name: "",
				},
			},
		},
	}

	w := new(bytes.Buffer)
	c := New(config)
	c.output = w

	result := c.Run()
	must.One(t, result)
	message := strings.TrimSpace(w.String())
	must.Eq(t, `babycli: component name missing`, message)
}

func TestComponent_validate_name_single(t *testing.T) {
	t.Parallel()

	config := &Configuration{
		Top: &Component{
			Components: Components{
				{
					Name: "first",
				},
				{
					Name: "x",
				},
			},
		},
	}

	w := new(bytes.Buffer)
	c := New(config)
	c.output = w

	result := c.Run()
	must.One(t, result)
	message := strings.TrimSpace(w.String())
	must.Eq(t, `babycli: component "x" must be more than one character`, message)
}
