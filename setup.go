// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"context"
	"io"
	"math"
	"os"
	"slices"

	"noxide.lol/go/stacks"
)

type Code = int

const (
	Success   Code = 0
	Failure   Code = 1
	Usability Code = math.MaxInt
)

type result struct {
	code Code

	// TODO: should this be used or removed?
	// message string
}

type Configuration struct {
	Arguments []string
	Top       *Component
	Globals   Flags
	Version   string
	Output    io.Writer
	Context   context.Context
}

func Arguments() []string {
	return os.Args[1:]
}

func New(c *Configuration) *Runnable {
	arguments := slices.Clone(c.Arguments)
	slices.Reverse(arguments)
	c.Top.args = stacks.Simple(arguments...)
	c.Top.version = c.Version
	c.Top.globals = c.globals()
	c.Top.context = c.context()
	output := c.Output
	if output == nil {
		output = os.Stderr
	}
	return &Runnable{
		root:   c.Top,
		output: output,
	}
}

func (c *Configuration) context() context.Context {
	if c.Context == nil {
		return context.Background()
	}
	return c.Context
}

func (c *Configuration) globals() Flags {
	return append(c.Globals, helpFlag)
}

type Runnable struct {
	root   *Component
	output io.Writer
}

func (r *Runnable) Run() (c Code) {
	defer func() {
		if p := recover(); p != nil {
			_, _ = io.WriteString(r.output, p.(string))
			c = Failure
		}
	}()
	result := r.run()
	return result.code
}

func (r *Runnable) run() *result {
	return r.root.run(r.output)
}
