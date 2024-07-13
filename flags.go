// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type FlagType uint8

const (
	StringFlag FlagType = iota
	IntFlag
	BoolFlag
	DurationFlag
)

func (t FlagType) String() string {
	switch t {
	case StringFlag:
		return "string"
	case IntFlag:
		return "integer"
	case BoolFlag:
		return "boolean"
	case DurationFlag:
		return "duration"
	}
	panic("babycli: not a flag type")
}

type Flag struct {
	Type    FlagType
	Require bool
	Repeats bool
	Long    string
	Short   string
	Usage   string
}

func (f *Flag) help() [3]string {
	var parts [3]string
	switch {
	case f.Long != "" && f.Short != "":
		parts[0] = fmt.Sprintf("--%s/-%s", f.Long, f.Short)
	case f.Long != "":
		parts[0] = "--" + f.Long
	default:
		parts[0] = "-" + f.Short
	}
	parts[1] = f.Type.String()
	parts[2] = f.Usage
	return parts
}

func (f *Flag) Identity() string {
	if f.Long == "" {
		return f.Short
	}
	return f.Long
}

func (f *Flag) Is(name string) bool {
	if len(name) == 1 {
		return f.Short == name
	}
	return f.Long == name
}

type Flags []*Flag

func (fs Flags) Contains(name string) bool {
	return slices.ContainsFunc(fs, func(f *Flag) bool {
		return f.Is(name)
	})
}

func (fs Flags) Get(name string) *Flag {
	for _, f := range fs {
		if f.Is(name) {
			return f
		}
	}
	panicf("flag %q is not defined", name)
	return nil
}

func (fs Flags) write(w io.Writer) {
	lines := make([][3]string, 0, len(fs))
	for _, flag := range fs {
		lines = append(lines, flag.help())
	}

	var max0, max1 int

	for i := 0; i < len(lines); i++ {
		max0 = max(max0, len(lines[i][0]))
		max1 = max(max1, len(lines[i][1]))
	}

	for _, line := range lines {
		_, _ = io.WriteString(w, rightPad(max0, line[0]))
		_, _ = io.WriteString(w, leftPad(max1, line[1]))
		_, _ = io.WriteString(w, "- ")
		_, _ = io.WriteString(w, line[2])
		_, _ = io.WriteString(w, "\n")
	}
}

func leftPad(size int, s string) string {
	sb := new(strings.Builder)
	n := (size + 1) - len(s)
	for i := 0; i < n; i++ {
		sb.WriteString(" ")
	}
	sb.WriteString(s)
	sb.WriteString(" ")
	return sb.String()
}

func rightPad(size int, s string) string {
	sb := new(strings.Builder)
	sb.WriteString(s)
	n := (size + 1) - len(s)
	for i := 0; i < n; i++ {
		sb.WriteString(" ")
	}
	return sb.String()
}
