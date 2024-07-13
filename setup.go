// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"slices"

	"noxide.lol/go/stacks"
)

type ExitCode uint8

type Result struct {
	Code    ExitCode
	Message string
}

const (
	ExitSuccess ExitCode = iota
	ExitFailure
)

type Configuration struct {
	root *Component
}

func New(args []string, root *Component) *Configuration {
	arguments := slices.Clone(args)
	slices.Reverse(arguments)
	root.args = stacks.Simple(arguments...)
	return &Configuration{
		root: root,
	}
}

func (r *Configuration) Run() ExitCode {
	if r := recover(); r != nil {
		return ExitFailure
	}
	result := r.run()
	return result.Code
}

func (r *Configuration) run() *Result {
	return r.root.run()
}
