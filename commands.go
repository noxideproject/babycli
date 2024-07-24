// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"context"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"

	"noxide.lol/go/stacks"
)

type Func func(*Component) Code

type values struct {
	strings   map[string][]string
	ints      map[string][]int
	bools     map[string][]bool
	durations map[string][]time.Duration
}

func (v *values) stringCount(flag string) int {
	return len(v.strings[flag])
}

func (v *values) intCount(flag string) int {
	return len(v.ints[flag])
}

func (v *values) boolCount(flag string) int {
	return len(v.bools[flag])
}

func (v *values) durationCount(flag string) int {
	return len(v.durations[flag])
}

func (v *values) helpSet() bool {
	for k, bs := range v.bools {
		if k == "help" || k == "h" {
			for _, b := range bs {
				if b {
					return true
				}
			}
		}
	}
	return false
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

type Component struct {
	Name string

	Help string

	Description string

	Components Components

	Function Func

	Flags Flags

	args stacks.Stack[string]

	flat []string

	vals *values

	globals Flags

	version string

	context context.Context
}

func (c *Component) Context() context.Context {
	return c.context
}

func (c *Component) Arguments() []string {
	if len(c.flat) == 0 && c.args.Size() > 0 {
		c.flat = make([]string, 0, c.args.Size())
		for i := 0; i < c.args.Size(); i++ {
			c.flat = append(c.flat, c.args.Pop())
		}
	}
	return c.flat
}

func (c *Component) Nargs() int {
	return len(c.Arguments())
}

func (c *Component) Leaf() bool {
	return len(c.Components) == 0
}

func (c *Component) init() {
	if c.vals == nil {
		c.vals = &values{
			strings:   make(map[string][]string, 0),
			ints:      make(map[string][]int, 0),
			bools:     make(map[string][]bool, 0),
			durations: make(map[string][]time.Duration, 0),
		}
	}
}

func (c *Component) run(output io.Writer) *result {
	c.init()

	for !c.args.Empty() {
		if more := c.processFlags(); !more {
			break
		}
	}

	if c.vals.helpSet() {
		text := c.help()
		writef(output, text)
		return &result{code: Success}
	}

	if c.Leaf() && c.Function != nil {
		code := c.Function(c)
		if code == Usability {
			text := c.help()
			writef(output, text)
			return &result{code: Failure}
		}
		return &result{code: code}
	}

	if c.args.Empty() {
		text := c.help()
		writef(output, text)
		return &result{code: Failure}
	}

	sub := c.args.Pop()
	cmd := c.Components.Get(sub)
	cmd.args = c.args
	cmd.vals = c.vals
	cmd.globals = c.globals
	cmd.context = c.context
	return cmd.run(output)
}

func (c *Component) processFlags() bool {
	arg := c.args.Peek()

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
	combine := make(Flags, 0, len(c.Flags)+len(c.globals))
	combine = append(combine, c.Flags...)
	combine = append(combine, c.globals...)

	name := c.args.Pop()
	name = strings.TrimLeft(name, "-")
	flag := combine.Get(name)

	switch flag.Type {
	case BooleanFlag:
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

func (c *Component) HasString(flag string) bool {
	return c.vals.stringCount(flag) > 0
}

func (c *Component) GetString(flag string) string {
	switch c.vals.stringCount(flag) {
	case 0:
		f := c.Flags.Get(flag)
		if f.Default != nil {
			return f.Default.Value.(string)
		}
		if f.Require {
			panicf("no value for string flag %q", flag)
		}
	case 1:
		return c.vals.strings[flag][0]
	default:
		panicf("multiple values set for string flag %q", flag)
	}
	return ""
}

func (c *Component) GetStrings(flag string) []string {
	if n := c.vals.stringCount(flag); n == 0 {
		f := c.Flags.Get(flag)
		if f.Default != nil {
			return []string{f.Default.Value.(string)}
		}
		if f.Require {
			panicf("no value for string flag %q", flag)
		}
	}
	return slices.Clone(c.vals.strings[flag])
}

func (c *Component) HasInt(flag string) bool {
	return c.vals.intCount(flag) > 0
}

func (c *Component) GetInt(flag string) int {
	switch c.vals.intCount(flag) {
	case 0:
		f := c.Flags.Get(flag)
		if f.Default != nil {
			return f.Default.Value.(int)
		}
		if f.Require {
			panicf("no value for int flag %q", flag)
		}
	case 1:
		return c.vals.ints[flag][0]
	default:
		panicf("multiple values set for int flag %q", flag)
	}
	return 0
}

func (c *Component) GetInts(flag string) []int {
	if n := c.vals.intCount(flag); n == 0 {
		f := c.Flags.Get(flag)
		if f.Default != nil {
			return []int{f.Default.Value.(int)}
		}
		if f.Require {
			panicf("no value for int flag %q", flag)
		}
	}
	return slices.Clone(c.vals.ints[flag])
}

func (c *Component) HasDuration(flag string) bool {
	return c.vals.durationCount(flag) > 0
}

func (c *Component) GetDuration(flag string) time.Duration {
	switch c.vals.durationCount(flag) {
	case 0:
		f := c.Flags.Get(flag)
		if f.Default != nil {
			return f.Default.Value.(time.Duration)
		}
		if f.Require {
			panicf("no value for duration flag %q", flag)
		}
	case 1:
		return c.vals.durations[flag][0]
	default:
		panicf("multiple values set for duration flag %q", flag)
	}
	return 0
}

func (c *Component) GetDurations(flag string) []time.Duration {
	if n := c.vals.intCount(flag); n == 0 {
		f := c.Flags.Get(flag)
		if f.Default != nil {
			return []time.Duration{f.Default.Value.(time.Duration)}
		}
		if f.Require {
			panicf("no value for duration flag %q", flag)
		}
	}
	return slices.Clone(c.vals.durations[flag])
}

func (c *Component) HasBool(flag string) bool {
	return c.vals.boolCount(flag) > 0
}

func (c *Component) GetBool(flag string) bool {
	switch c.vals.boolCount(flag) {
	case 0:
		f := c.Flags.Get(flag)
		if f.Default != nil {
			return f.Default.Value.(bool)
		}
		if f.Require {
			panicf("no value for boolean flag %q", flag)
		}
	case 1:
		return c.vals.bools[flag][0]
	default:
		panicf("multiple values set for boolean flag %q", flag)
	}
	return false
}

func (c *Component) GetBools(flag string) []bool {
	if n := c.vals.boolCount(flag); n == 0 {
		f := c.Flags.Get(flag)
		if f.Default != nil {
			return []bool{f.Default.Value.(bool)}
		}
		if f.Require {
			panicf("no value for boolean flag %q", flag)
		}
	}
	return slices.Clone(c.vals.bools[flag])
}
