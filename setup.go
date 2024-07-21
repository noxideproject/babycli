// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"context"
	"io"
	"os"
	"slices"

	"noxide.lol/go/stacks"
)

type ExitCode = int

type Result struct {
	Code    ExitCode
	Message string
}

const (
	ExitSuccess ExitCode = iota
	ExitFailure
)

type Configuration struct {
	Arguments []string
	Top       *Component
	Globals   Flags
	Version   string
	Output    io.Writer
	Context   context.Context
}

type Runnable struct {
	root   *Component
	output io.Writer
}

func Arguments() []string {
	return os.Args[1:]
}

func New(c *Configuration) *Runnable {
	arguments := slices.Clone(c.Arguments)
	slices.Reverse(arguments)
	c.Top.args = stacks.Simple(arguments...)
	c.Top.version = c.Version
	c.Top.globals = append(c.Globals, helpFlag)
	c.Top.Context = c.context()
	output := c.Output
	if output == nil {
		output = os.Stderr
	}
	return &Runnable{
		root:   c.Top,
		output: output,
	}
}

func (r *Runnable) Run() ExitCode {
	if r := recover(); r != nil {
		return ExitFailure
	}
	result := r.run()
	return result.Code
}

func (r *Runnable) run() *Result {
	return r.root.run(r.output)
}

func (c *Configuration) context() context.Context {
	if c.Context == nil {
		return context.Background()
	}
	return c.Context
}
