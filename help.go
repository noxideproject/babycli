// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"strings"
)

var helpFlag = &Flag{
	Type:    BoolFlag,
	Require: false,
	Repeats: false,
	Long:    "help",
	Short:   "h",
	Usage:   "print help message",
}

const (
	tab = "  "
)

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
		for _, cmp := range c.Components {
			sb.WriteString(tab)
			sb.WriteString(cmp.Name)
			sb.WriteString(" - ")
			sb.WriteString(cmp.Help)
			sb.WriteString("\n")
		}
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

	return sb.String()
}

func chop(s string) []string {
	s = strings.TrimSpace(s)
	return strings.Split(s, "\n")
}
