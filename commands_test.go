// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

type testCase struct {
	name string
	args []string
	root *Component
	exp  string
}

func TestRun_topCommand(t *testing.T) {
	t.Parallel()

	var output string

	cases := []testCase{
		{
			name: "bare",
			exp:  "ok",
			args: nil,
			root: &Component{
				Function: func(*Component) {
					output = "ok"
				},
			},
		},
		{
			name: "string flag long",
			exp:  "hello, bob!",
			args: []string{"--name", "bob"},
			root: &Component{
				Flags: Flags{
					{
						Type:  StringFlag,
						Long:  "name",
						Short: "n",
					},
				},
				Function: func(c *Component) {
					name := c.GetString("name")
					output = fmt.Sprintf("hello, %s!", name)
				},
			},
		},
		{
			name: "string flag short",
			exp:  "hello, carol!",
			args: []string{"-n", "carol"},
			root: &Component{
				Flags: Flags{
					{
						Type:  StringFlag,
						Long:  "name",
						Short: "n",
					},
				},
				Function: func(c *Component) {
					name := c.GetString("name")
					output = fmt.Sprintf("hello, %s!", name)
				},
			},
		},
		{
			name: "string flags multi",
			exp:  "hello alice bob carol",
			args: []string{"--name", "alice", "-n", "bob", "--name", "carol"},
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Short:   "n",
						Repeats: true,
					},
				},
				Function: func(c *Component) {
					names := c.GetStrings("name")
					output = "hello " + strings.Join(names, " ")
				},
			},
		},
		{
			name: "int flag long",
			exp:  "age is 34",
			args: []string{"--age", "34"},
			root: &Component{
				Flags: Flags{
					{
						Type: IntFlag,
						Long: "age",
					},
				},
				Function: func(c *Component) {
					age := c.GetInt("age")
					output = fmt.Sprintf("age is %d", age)
				},
			},
		},
		{
			name: "int flag short",
			exp:  "age is 30",
			args: []string{"-a", "30"},
			root: &Component{
				Flags: Flags{
					{
						Type:  IntFlag,
						Long:  "age",
						Short: "a",
					},
				},
				Function: func(c *Component) {
					age := c.GetInt("age")
					output = fmt.Sprintf("age is %d", age)
				},
			},
		},
		{
			name: "int flags multi",
			exp:  "ages are 20 30 40 50",
			args: []string{"-a", "20", "--age", "30", "--age", "40", "-a", "50"},
			root: &Component{
				Flags: Flags{
					{
						Type:  IntFlag,
						Long:  "age",
						Short: "a",
					},
				},
				Function: func(c *Component) {
					ages := c.GetInts("age")
					output = fmt.Sprintf("ages are %d %d %d %d", ages[0], ages[1], ages[2], ages[3])
				},
			},
		},
		{
			name: "duration flag long",
			exp:  "ttl is 2m0s",
			args: []string{"--ttl", "120s"},
			root: &Component{
				Flags: Flags{
					{
						Type: DurationFlag,
						Long: "ttl",
					},
				},
				Function: func(c *Component) {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("ttl is %s", ttl)
				},
			},
		},
		{
			name: "duration flag short",
			exp:  "ttl is 3m0s",
			args: []string{"-s", "180s"},
			root: &Component{
				Flags: Flags{
					{
						Type:  DurationFlag,
						Long:  "ttl",
						Short: "s",
					},
				},
				Function: func(c *Component) {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("ttl is %s", ttl)
				},
			},
		},
		{
			name: "multi duration flags",
			exp:  "ttls are 2m0s 3m0s",
			args: []string{"-s", "120s", "--ttl", "180s"},
			root: &Component{
				Flags: Flags{
					{
						Type:  DurationFlag,
						Long:  "ttl",
						Short: "s",
					},
				},
				Function: func(c *Component) {
					ttls := c.GetDurations("ttl")
					output = fmt.Sprintf("ttls are %s %s", ttls[0], ttls[1])
				},
			},
		},
		{
			name: "bool flag long explicit",
			exp:  "force is true",
			args: []string{"--force", "true"},
			root: &Component{
				Flags: Flags{
					{
						Type: BoolFlag,
						Long: "force",
					},
				},
				Function: func(c *Component) {
					f := c.GetBool("force")
					output = fmt.Sprintf("force is %t", f)
				},
			},
		},
		{
			name: "bool flag long implicit",
			exp:  "force is true",
			args: []string{"--force"},
			root: &Component{
				Flags: Flags{
					{
						Type: BoolFlag,
						Long: "force",
					},
				},
				Function: func(c *Component) {
					f := c.GetBool("force")
					output = fmt.Sprintf("force is %t", f)
				},
			},
		},
		{
			name: "bool flag long explicit false",
			exp:  "force is false",
			args: []string{"--force", "false"},
			root: &Component{
				Flags: Flags{
					{
						Type: BoolFlag,
						Long: "force",
					},
				},
				Function: func(c *Component) {
					f := c.GetBool("force")
					output = fmt.Sprintf("force is %t", f)
				},
			},
		},
		{
			name: "bool flag short",
			exp:  "force is true",
			args: []string{"-f"},
			root: &Component{
				Flags: Flags{
					{
						Type:  BoolFlag,
						Long:  "force",
						Short: "f",
					},
				},
				Function: func(c *Component) {
					f := c.GetBool("force")
					output = fmt.Sprintf("force is %t", f)
				},
			},
		},
	}

	for _, tc := range cases {
		output = "" // reset for each test case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.exp, output)
			must.Eq(t, ExitSuccess, result)
		})
	}
}

func TestRun_childCommand(t *testing.T) {
	t.Parallel()

	var output string
	cases := []testCase{
		{
			name: "bare",
			exp:  "this is about",
			args: []string{"about"},
			root: &Component{
				Components: Components{
					{
						Name: "about",
						Function: func(*Component) {
							output = "this is about"
						},
					},
				},
			},
		},
		{
			name: "string flag long",
			exp:  "hello, bob!",
			args: []string{"sayhi", "--name", "bob"},
			root: &Component{
				Components: Components{
					{
						Name: "sayhi",
						Flags: Flags{
							{
								Type: StringFlag,
								Long: "name",
							},
						},
						Function: func(c *Component) {
							name := c.GetString("name")
							output = fmt.Sprintf("hello, %s!", name)
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		output = "" // reset for each test case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.exp, output)
			must.Eq(t, ExitSuccess, result)
		})
	}
}

func TestRun_grandchildCommand(t *testing.T) {
	t.Parallel()

	var output string

	cases := []testCase{
		{
			name: "bare",
			exp:  "this is grandchild",
			args: []string{"first", "second"},
			root: &Component{
				Components: Components{
					{
						Name: "first",
						Components: Components{
							{
								Name: "second",
								Function: func(*Component) {
									output = "this is grandchild"
								},
							},
						},
					},
				},
			},
		},
		{
			name: "string flag long",
			exp:  "hello, carol!",
			args: []string{"greeting", "hello", "--name", "carol"},
			root: &Component{
				Components: Components{
					{
						Name: "greeting",
						Components: Components{
							{
								Name: "hello",
								Flags: Flags{
									{
										Type: StringFlag,
										Long: "name",
									},
								},
								Function: func(c *Component) {
									name := c.GetString("name")
									output = fmt.Sprintf("hello, %s!", name)
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		output = "" // reset for each test case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.exp, output)
			must.Eq(t, ExitSuccess, result)
		})
	}
}

func TestHelp_top(t *testing.T) {
	t.Parallel()

	config := &Configuration{
		Arguments: nil,
		Top: &Component{
			Name: "program",
		},
	}

	c := New(config)

	code := c.Run()
	must.Eq(t, ExitFailure, code)
}
