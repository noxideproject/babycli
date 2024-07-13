// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

// Package babycli is a library for declrative parsing of command line arguments,
// including support for sub-commands, command aliases, long and short flag names,
// repeated flags, and custom help messages.
package babycli

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"noxide.lol/go/stacks"
)

type Func func(*Component)

type FlagType uint8

const (
	StringFlag FlagType = iota
	IntFlag
	BoolFlag
	DurationFlag
)

type Flag struct {
	Type    FlagType
	Require bool
	Repeats bool
	Long    string
	Short   string
	Help    string
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

type values struct {
	strings   map[string][]string
	ints      map[string][]int
	bools     map[string][]bool
	durations map[string][]time.Duration
}

type Components []*Component

func (cs Components) Contains(name string) bool {
	return slices.ContainsFunc(cs, func(c *Component) bool {
		return c.Name == name
	})
}

func (cs Components) Get(name string) *Component {
	for _, c := range cs {
		if c.Name == name {
			return c
		}
	}
	panicf("subcommand %q is not defined", name)
	return nil
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

type Component struct {
	Name string

	Context context.Context

	Components Components

	Function Func

	Flags Flags

	args stacks.Stack[string]

	vals *values
}

func (c *Component) Leaf() bool {
	if len(c.Components) > 0 {
		return false
	}
	if c.Function == nil {
		panicf("no function for leaf command %q", c.Name)
	}
	return true
}

func (c *Component) run() *Result {
	if c.vals == nil {
		c.vals = &values{
			strings:   make(map[string][]string, 0),
			ints:      make(map[string][]int, 0),
			bools:     make(map[string][]bool, 0),
			durations: make(map[string][]time.Duration, 0),
		}
	}

	for !c.args.Empty() {
		if !c.setupFlags() {
			break
		}
	}

	if c.Leaf() {
		c.Function(c)
		return &Result{Code: ExitSuccess}
	}
	sub := c.args.Pop()
	cmd := c.Components.Get(sub)
	cmd.args = c.args
	cmd.vals = c.vals
	return cmd.run()
}

func (c *Component) setupFlags() bool {
	arg := c.args.Peek()
	fmt.Println("next() arg", arg)
	switch {
	case strings.HasPrefix(arg, "--"):
		c.consumeFlag()
		return true
	case strings.HasPrefix(arg, "-"):
		c.consumeFlag()
		return true
	default:
		return false
	}
}

func (c *Component) consumeFlag() {
	name := c.args.Pop()
	name = strings.TrimLeft(name, "-")
	flag := c.Flags.Get(name)

	switch flag.Type {
	case BoolFlag:
		c.consumeBoolFlag(flag.Identity())
	case StringFlag:
		c.consumeStringFlag(flag.Identity())
	case IntFlag:
		c.consumeIntFlag(flag.Identity())
	case DurationFlag:
		c.consumeDurationFlag(flag.Identity())
	}
}

func (c *Component) consumeBoolFlag(identity string) {
	if c.args.Empty() {
		c.vals.bools[identity] = append(c.vals.bools[identity], true)
		return
	}

	next := c.args.Peek()
	switch {
	case next == "true":
		c.vals.bools[identity] = append(c.vals.bools[identity], true)
		_ = c.args.Pop()
	case next == "false":
		c.vals.bools[identity] = append(c.vals.bools[identity], false)
		_ = c.args.Pop()
	default:
		c.vals.bools[identity] = append(c.vals.bools[identity], true)
	}
}

func (c *Component) consumeStringFlag(identity string) {

	if c.args.Empty() {
		// TODO what about default values
		panicf("no value for string flag %q", identity)
	}

	if strings.HasPrefix(c.args.Peek(), "-") {
		panicf("no value for string flag %q", identity)
	}

	value := c.args.Pop()
	c.vals.strings[identity] = append(c.vals.strings[identity], value)
}

func (c *Component) consumeIntFlag(identity string) {

	if c.args.Empty() {
		// TODO what about default values
		panicf("no value for int flag %q", identity)
	}

	if strings.HasPrefix(c.args.Peek(), "-") {
		panicf("no value for int flag %q", identity)
	}

	value := c.args.Pop()
	i, err := strconv.Atoi(value)
	if err != nil {
		panicf("unable to convert value for flag %q to int %q", identity, value)
	}
	c.vals.ints[identity] = append(c.vals.ints[identity], i)
}

func (c *Component) consumeDurationFlag(identity string) {
	if c.args.Empty() {
		// TODO what about default values
		panicf("no value for string flag %q", identity)
	}

	if strings.HasPrefix(c.args.Peek(), "-") {
		panicf("no value for string flag %q", identity)
	}

	value := c.args.Pop()
	dur, err := time.ParseDuration(value)
	if err != nil {
		panicf("unable to convert value for flag %q to duration %q", identity, value)
	}
	c.vals.durations[identity] = append(c.vals.durations[identity], dur)
}

// command -v

// command argument

// command subcommand -v

// command -f
// command -f true
// command -f=true

func (c *Component) GetString(flag string) string {
	if len(c.vals.strings[flag]) == 0 {
		panicf("no value for string flag %q", flag)
	}
	if len(c.vals.strings[flag]) > 1 {
		panicf("multiple values for string flag %q", flag)
	}
	return c.vals.strings[flag][0]
}

func (c *Component) GetStrings(flag string) []string {
	return c.vals.strings[flag]
}

func (c *Component) GetInt(flag string) int {
	if len(c.vals.ints[flag]) == 0 {
		panicf("no value for int flag %q", flag)
	}
	if len(c.vals.ints[flag]) > 1 {
		panicf("multiple values for int flag %q", flag)
	}
	return c.vals.ints[flag][0]
}

func (c *Component) GetInts(flag string) []int {
	return c.vals.ints[flag]
}

func (c *Component) GetDuration(flag string) time.Duration {
	if len(c.vals.durations[flag]) == 0 {
		panicf("no value for duration flag %q", flag)
	}
	if len(c.vals.durations[flag]) > 1 {
		panicf("multiple values for duration flag %q", flag)
	}
	return c.vals.durations[flag][0]
}

func (c *Component) GetDurations(flag string) []time.Duration {
	return c.vals.durations[flag]
}

func (c *Component) GetBool(flag string) bool {
	if len(c.vals.bools[flag]) == 0 {
		panicf("no value for bool flag %q", flag)
	}
	if len(c.vals.bools[flag]) > 1 {
		panicf("multiple values for bool flag %q", flag)
	}
	return c.vals.bools[flag][0]
}
