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
				Function: func(*Component) Code {
					output = "ok"
					return Success
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
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = fmt.Sprintf("hello, %s!", name)
					return Success
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
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = fmt.Sprintf("hello, %s!", name)
					return Success
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
				Function: func(c *Component) Code {
					names := c.GetStrings("name")
					output = "hello " + strings.Join(names, " ")
					return Success
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
				Function: func(c *Component) Code {
					age := c.GetInt("age")
					output = fmt.Sprintf("age is %d", age)
					return Success
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
				Function: func(c *Component) Code {
					age := c.GetInt("age")
					output = fmt.Sprintf("age is %d", age)
					return Success
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
				Function: func(c *Component) Code {
					ages := c.GetInts("age")
					output = fmt.Sprintf("ages are %d %d %d %d", ages[0], ages[1], ages[2], ages[3])
					return Success
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
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("ttl is %s", ttl)
					return Success
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
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("ttl is %s", ttl)
					return Success
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
				Function: func(c *Component) Code {
					ttls := c.GetDurations("ttl")
					output = fmt.Sprintf("ttls are %s %s", ttls[0], ttls[1])
					return Success
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
				Function: func(c *Component) Code {
					f := c.GetBool("force")
					output = fmt.Sprintf("force is %t", f)
					return Success
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
				Function: func(c *Component) Code {
					f := c.GetBool("force")
					output = fmt.Sprintf("force is %t", f)
					return Success
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
				Function: func(c *Component) Code {
					f := c.GetBool("force")
					output = fmt.Sprintf("force is %t", f)
					return Success
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
				Function: func(c *Component) Code {
					f := c.GetBool("force")
					output = fmt.Sprintf("force is %t", f)
					return Success
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
			must.Eq(t, Success, result)
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
						Function: func(*Component) Code {
							output = "this is about"
							return Success
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
						Function: func(c *Component) Code {
							name := c.GetString("name")
							output = fmt.Sprintf("hello, %s!", name)
							return Success
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
			must.Eq(t, Success, result)
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
								Function: func(*Component) Code {
									output = "this is grandchild"
									return Success
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
								Function: func(c *Component) Code {
									name := c.GetString("name")
									output = fmt.Sprintf("hello, %s!", name)
									return Success
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
			must.Eq(t, Success, result)
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
	must.Eq(t, Failure, code)
}
