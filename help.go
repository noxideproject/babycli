// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"io"
	"strings"
)

var helpFlag = &Flag{
	Type:    BooleanFlag,
	Require: false,
	Repeats: false,
	Long:    "help",
	Short:   "h",
	Help:    "print help message",
}

const (
	tab = "  "
)

func (c Components) write(w io.Writer) {
	lines := make([][2]string, 0, len(c))

	for _, component := range c {
		lines = append(lines, [2]string{component.Name, component.Help})
	}

	var max0 int

	for i := 0; i < len(lines); i++ {
		max0 = max(max0, len(lines[i][0]))
	}

	for _, line := range lines {
		_, _ = io.WriteString(w, "  ")
		_, _ = io.WriteString(w, rightPad(max0, line[0]))
		_, _ = io.WriteString(w, "- ")
		_, _ = io.WriteString(w, line[1])
		_, _ = io.WriteString(w, "\n")
	}
}

func (c *Component) help() string {
	sb := new(strings.Builder)
	sb.WriteString("NAME:\n")
	sb.WriteString(tab)
	sb.WriteString(c.Name)
	if c.Help != "" {
		sb.WriteString(" - ")
		sb.WriteString(c.Help)
	}
	sb.WriteString("\n\n")

	sb.WriteString("USAGE:\n")
	sb.WriteString(tab)
	sb.WriteString(c.Name)
	sb.WriteString(tab)
	sb.WriteString("[global options] [command [command options]] [arguments...]")
	sb.WriteString("\n\n")

	if c.version != "" {
		sb.WriteString("VERSION:\n")
		sb.WriteString(tab)
		sb.WriteString(c.version)
		sb.WriteString("\n\n")
	}

	if c.Description != "" {
		sb.WriteString("DESCRIPTION:\n")
		lines := chop(c.Description)
		for _, line := range lines {
			sb.WriteString(tab)
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	if len(c.Components) > 0 {
		sb.WriteString("COMMANDS:\n")
		c.Components.write(sb)
		sb.WriteString("\n")
	}

	if len(c.Flags) > 0 {
		sb.WriteString("OPTIONS:\n")
		c.Flags.write(sb)
		sb.WriteString("\n")
	}

	if len(c.globals) > 0 {
		sb.WriteString("GLOBALS:\n")
		c.globals.write(sb)
		sb.WriteString("\n")
	}

	s := sb.String()
	return strings.TrimSpace(s)
}

func chop(s string) []string {
	s = strings.TrimSpace(s)
	return strings.Split(s, "\n")
}
