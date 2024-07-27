// Copyright (c) The Noxide Project Authors
// SPDX-License-Identifier: BSD-3-Clause

package babycli

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	"noxide.lol/go/stacks"
)

type testCase struct {
	name     string
	args     []string
	root     *Component
	expText  string
	expCode  Code
	expPanic string
}

func TestRun_topCommand(t *testing.T) {
	t.Parallel()

	var output string

	cases := []testCase{
		{
			name:    "bare",
			expText: "ok",
			args:    nil,
			root: &Component{
				Function: func(*Component) Code {
					output = "ok"
					return Success
				},
			},
		},
		{
			name:    "string flag long",
			expText: "hello, bob!",
			args:    []string{"--name", "bob"},
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
			name:    "string flag short",
			expText: "hello, carol!",
			args:    []string{"-n", "carol"},
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
			name:    "string flags multi",
			expText: "hello alice bob carol",
			args:    []string{"--name", "alice", "-n", "bob", "--name", "carol"},
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
			name:    "int flag long",
			expText: "age is 34",
			args:    []string{"--age", "34"},
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
			name:    "int flag short",
			expText: "age is 30",
			args:    []string{"-a", "30"},
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
			name:    "int flags multi",
			expText: "ages are 20 30 40 50",
			args:    []string{"-a", "20", "--age", "30", "--age", "40", "-a", "50"},
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
			name:    "duration flag long",
			expText: "ttl is 2m0s",
			args:    []string{"--ttl", "120s"},
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
			name:    "duration flag short",
			expText: "ttl is 3m0s",
			args:    []string{"-s", "180s"},
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
			name:    "multi duration flags",
			expText: "ttls are 2m0s 3m0s",
			args:    []string{"-s", "120s", "--ttl", "180s"},
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
			name:    "bool flag long explicit",
			expText: "force is true",
			args:    []string{"--force", "true"},
			root: &Component{
				Flags: Flags{
					{
						Type: BooleanFlag,
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
			name:    "bool flag long implicit",
			expText: "force is true",
			args:    []string{"--force"},
			root: &Component{
				Flags: Flags{
					{
						Type: BooleanFlag,
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
			name:    "bool flag long explicit false",
			expText: "force is false",
			args:    []string{"--force", "false"},
			root: &Component{
				Flags: Flags{
					{
						Type: BooleanFlag,
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
			name:    "bool flag short",
			expText: "force is true",
			args:    []string{"-f"},
			root: &Component{
				Flags: Flags{
					{
						Type:  BooleanFlag,
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
			must.Eq(t, tc.expText, output)
			must.Eq(t, Success, result)
		})
	}
}

func TestRun_childCommand(t *testing.T) {
	t.Parallel()

	var output string
	cases := []testCase{
		{
			name:    "bare",
			expText: "this is about",
			args:    []string{"about"},
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
			name:    "string flag long",
			expText: "hello, bob!",
			args:    []string{"sayhi", "--name", "bob"},
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
			must.Eq(t, tc.expText, output)
			must.Eq(t, Success, result)
		})
	}
}

func TestRun_grandchildCommand(t *testing.T) {
	t.Parallel()

	var output string

	cases := []testCase{
		{
			name:    "bare",
			expText: "this is grandchild",
			args:    []string{"first", "second"},
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
			name:    "string flag long",
			expText: "hello, carol!",
			args:    []string{"greeting", "hello", "--name", "carol"},
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
			must.Eq(t, tc.expText, output)
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

func TestComponent_GetString(t *testing.T) {
	t.Parallel()

	var output string
	var failure *strings.Builder

	cases := []testCase{
		{
			name:    "required string provided no default",
			expText: "hello bob",
			expCode: Success,
			args:    []string{"--name", "bob"},
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = "hello " + name
					return Success
				},
			},
		},
		{
			name:    "required string provided with default",
			expText: "hello bob",
			expCode: Success,
			args:    []string{"--name", "bob"},
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Require: true,
						Default: &Default{
							Value: "alice",
						},
					},
				},
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = "hello " + name
					return Success
				},
			},
		},
		{
			name:    "required string not provided with default",
			expText: "hello alice",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Require: true,
						Default: &Default{
							Value: "alice",
						},
					},
				},
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = "hello " + name
					return Success
				},
			},
		},
		{
			name:     "required string not provided no default",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: no value for string flag "name"`,
			args:     nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = "hello " + name
					return Success
				},
			},
		},
		{
			name:    "optional string not provided no default",
			expText: "hello .",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = fmt.Sprintf("hello %s.", name)
					return Success
				},
			},
		},
		{
			name:     "no repeat string provided twice",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: multiple values set for string flag "name"`,
			args:     []string{"--name", "bob", "--name", "carl"},
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Require: false,
						Repeats: false,
					},
				},
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = fmt.Sprintf("hello %s.", name)
					return Success
				},
			},
		},
		{
			name:     "repeat string provided twice",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: multiple values set for string flag "name"`,
			args:     []string{"--name", "bob", "--name", "carl"},
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Require: false,
						Repeats: true,
					},
				},
				Function: func(c *Component) Code {
					name := c.GetString("name") // must use GetStrings
					output = fmt.Sprintf("hello %s.", name)
					return Success
				},
			},
		},
		{
			name:    "use equal sign",
			expText: "hello bob.",
			expCode: Success,
			args:    []string{"--name=bob"},
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Require: false,
						Repeats: false,
					},
				},
				Function: func(c *Component) Code {
					name := c.GetString("name")
					output = fmt.Sprintf("hello %s.", name)
					return Success
				},
			},
		},
	}

	for _, tc := range cases {
		output = ""                    // reset for each case
		failure = new(strings.Builder) // reset for each case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
				Output:    failure,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.expText, output)
			must.Eq(t, tc.expCode, result)
			must.Eq(t, tc.expPanic, failure.String())
		})
	}
}

func TestComponent_GetStrings(t *testing.T) {
	t.Parallel()

	var output string
	var failure *strings.Builder

	cases := []testCase{
		{
			name:    "repeated strings provided no default",
			expText: "hello alice bob carl",
			expCode: Success,
			args:    []string{"--name", "alice", "--name", "bob", "--name", "carl"},
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Repeats: true,
						Require: false,
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
			name:    "repeated strings not provided no default not required",
			expText: "hello ",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Repeats: true,
						Require: false,
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
			name:     "repeated strings not provided no default required",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: no value for string flag "name"`,
			args:     nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Repeats: true,
						Require: true,
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
			name:    "repeated strings not provided with default required",
			expText: "hello dave",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "name",
						Repeats: true,
						Require: true,
						Default: &Default{
							Value: "dave",
						},
					},
				},
				Function: func(c *Component) Code {
					names := c.GetStrings("name")
					output = "hello " + strings.Join(names, " ")
					return Success
				},
			},
		},
	}

	for _, tc := range cases {
		output = ""                    // reset for each case
		failure = new(strings.Builder) // reset for each case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
				Output:    failure,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.expText, output)
			must.Eq(t, tc.expCode, result)
			must.Eq(t, tc.expPanic, failure.String())
		})
	}
}

func TestComponent_GetInt(t *testing.T) {
	t.Parallel()

	var output string
	var failure *strings.Builder

	cases := []testCase{
		{
			name:    "required int provided no default",
			expText: "hello 1",
			expCode: Success,
			args:    []string{"--age", "1"},
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					age := c.GetInt("age")
					output = fmt.Sprintf("hello %d", age)
					return Success
				},
			},
		},
		{
			name:    "required int provided with default",
			expText: "hello 2",
			expCode: Success,
			args:    []string{"--age", "2"},
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Require: true,
						Default: &Default{
							Value: 2,
						},
					},
				},
				Function: func(c *Component) Code {
					age := c.GetInt("age")
					output = fmt.Sprintf("hello %d", age)
					return Success
				},
			},
		},
		{
			name:    "required int not provided with default",
			expText: "hello 3",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    StringFlag,
						Long:    "age",
						Require: true,
						Default: &Default{
							Value: 3,
						},
					},
				},
				Function: func(c *Component) Code {
					age := c.GetInt("age")
					output = fmt.Sprintf("hello %d", age)
					return Success
				},
			},
		},
		{
			name:     "required int not provided no default",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: no value for int flag "age"`,
			args:     nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					age := c.GetInt("age")
					output = fmt.Sprintf("hello %d", age)
					return Success
				},
			},
		},
		{
			name:    "optional int not provided no default",
			expText: "hello 0.",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					age := c.GetInt("age")
					output = fmt.Sprintf("hello %d.", age)
					return Success
				},
			},
		},
		{
			name:     "no repeat int provided twice",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: multiple values set for int flag "age"`,
			args:     []string{"--age", "4", "--age", "5"},
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Require: false,
						Repeats: false,
					},
				},
				Function: func(c *Component) Code {
					age := c.GetInt("age")
					output = fmt.Sprintf("hello %d.", age)
					return Success
				},
			},
		},
		{
			name:     "repeat int provided twice",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: multiple values set for int flag "age"`,
			args:     []string{"--age", "6", "--age", "7"},
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Require: false,
						Repeats: true,
					},
				},
				Function: func(c *Component) Code {
					age := c.GetInt("age") // must use GetInts
					output = fmt.Sprintf("hello %d.", age)
					return Success
				},
			},
		},
	}

	for _, tc := range cases {
		output = ""                    // reset for each case
		failure = new(strings.Builder) // reset for each case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
				Output:    failure,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.expText, output)
			must.Eq(t, tc.expCode, result)
			must.Eq(t, tc.expPanic, failure.String())
		})
	}
}

func TestComponent_GetInts(t *testing.T) {
	t.Parallel()

	var output string
	var failure *strings.Builder

	cases := []testCase{
		{
			name:    "repeated ints provided no default",
			expText: "hello [1 2 3]",
			expCode: Success,
			args:    []string{"--age", "1", "--age", "2", "--age", "3"},
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Repeats: true,
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					ages := c.GetInts("age")
					output = fmt.Sprintf("hello %v", ages)
					return Success
				},
			},
		},
		{
			name:    "repeated ints not provided no default not required",
			expText: "hello []",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Repeats: true,
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					ages := c.GetInts("age")
					output = fmt.Sprintf("hello %v", ages)
					return Success
				},
			},
		},
		{
			name:     "repeated ints not provided no default required",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: no value for int flag "age"`,
			args:     nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Repeats: true,
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					ages := c.GetInts("age")
					output = fmt.Sprintf("hello %v", ages)
					return Success
				},
			},
		},
		{
			name:    "repeated ints not provided with default required",
			expText: "hello [9]",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    IntFlag,
						Long:    "age",
						Repeats: true,
						Require: true,
						Default: &Default{
							Value: 9,
						},
					},
				},
				Function: func(c *Component) Code {
					ages := c.GetInts("age")
					output = fmt.Sprintf("hello %v", ages)
					return Success
				},
			},
		},
	}

	for _, tc := range cases {
		output = ""                    // reset for each case
		failure = new(strings.Builder) // reset for each case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
				Output:    failure,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.expText, output)
			must.Eq(t, tc.expCode, result)
			must.Eq(t, tc.expPanic, failure.String())
		})
	}
}

func TestComponent_GetDuration(t *testing.T) {
	t.Parallel()

	var output string
	var failure *strings.Builder

	cases := []testCase{
		{
			name:    "required duration provided no default",
			expText: "hello 1m0s",
			expCode: Success,
			args:    []string{"--ttl", "1m"},
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("hello %s", ttl)
					return Success
				},
			},
		},
		{
			name:    "required duration provided with default",
			expText: "hello 2m0s",
			expCode: Success,
			args:    []string{"--ttl", "2m0s"},
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Require: true,
						Default: &Default{
							Value: 3,
						},
					},
				},
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("hello %s", ttl)
					return Success
				},
			},
		},
		{
			name:    "required duration not provided with default",
			expText: "hello 4m0s",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Require: true,
						Default: &Default{
							Value: 4 * time.Minute,
						},
					},
				},
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("hello %s", ttl)
					return Success
				},
			},
		},
		{
			name:     "required duration not provided no default",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: no value for duration flag "ttl"`,
			args:     nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("hello %s", ttl)
					return Success
				},
			},
		},
		{
			name:    "optional duration not provided no default",
			expText: "hello 0s.",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("hello %s.", ttl)
					return Success
				},
			},
		},
		{
			name:     "no repeat duration provided twice",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: multiple values set for duration flag "ttl"`,
			args:     []string{"--ttl", "5m", "--ttl", "6m"},
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Require: false,
						Repeats: false,
					},
				},
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl")
					output = fmt.Sprintf("hello %s.", ttl)
					return Success
				},
			},
		},
		{
			name:     "repeat int provided twice",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: multiple values set for duration flag "ttl"`,
			args:     []string{"--ttl", "6m", "--ttl", "7m"},
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Require: false,
						Repeats: true,
					},
				},
				Function: func(c *Component) Code {
					ttl := c.GetDuration("ttl") // must use GetDurations
					output = fmt.Sprintf("hello %d.", ttl)
					return Success
				},
			},
		},
	}

	for _, tc := range cases {
		output = ""                    // reset for each case
		failure = new(strings.Builder) // reset for each case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
				Output:    failure,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.expText, output)
			must.Eq(t, tc.expCode, result)
			must.Eq(t, tc.expPanic, failure.String())
		})
	}
}

func TestComponent_GetDurations(t *testing.T) {
	t.Parallel()

	var output string
	var failure *strings.Builder

	cases := []testCase{
		{
			name:    "repeated durations provided no default",
			expText: "hello [1m0s 2m0s 3m0s]",
			expCode: Success,
			args:    []string{"--ttl", "1m", "--ttl", "2m", "--ttl", "3m"},
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Repeats: true,
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					ttls := c.GetDurations("ttl")
					output = fmt.Sprintf("hello %v", ttls)
					return Success
				},
			},
		},
		{
			name:    "repeated durations not provided no default not required",
			expText: "hello []",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Repeats: true,
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					ttls := c.GetDurations("ttl")
					output = fmt.Sprintf("hello %v", ttls)
					return Success
				},
			},
		},
		{
			name:     "repeated durations not provided no default required",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: no value for duration flag "ttl"`,
			args:     nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Repeats: true,
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					ttls := c.GetDurations("ttl")
					output = fmt.Sprintf("hello %v", ttls)
					return Success
				},
			},
		},
		{
			name:    "repeated durations not provided with default required",
			expText: "hello [9m0s]",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    DurationFlag,
						Long:    "ttl",
						Repeats: true,
						Require: true,
						Default: &Default{
							Value: 9 * time.Minute,
						},
					},
				},
				Function: func(c *Component) Code {
					ttls := c.GetDurations("ttl")
					output = fmt.Sprintf("hello %v", ttls)
					return Success
				},
			},
		},
	}

	for _, tc := range cases {
		output = ""                    // reset for each case
		failure = new(strings.Builder) // reset for each case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
				Output:    failure,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.expText, output)
			must.Eq(t, tc.expCode, result)
			must.Eq(t, tc.expPanic, failure.String())
		})
	}
}

func TestComponent_GetBoolean(t *testing.T) {
	t.Parallel()

	var output string
	var failure *strings.Builder

	cases := []testCase{
		{
			name:    "required boolean provided no default",
			expText: "hello true",
			expCode: Success,
			args:    []string{"--verbose", "true"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:    "required boolean provided with default",
			expText: "hello true",
			expCode: Success,
			args:    []string{"--verbose", "true"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: true,
						Default: &Default{
							Value: true,
						},
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:    "required boolean provided implicit no default",
			expText: "hello true",
			expCode: Success,
			args:    []string{"--verbose"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:    "required boolean provided false no default",
			expText: "hello false",
			expCode: Success,
			args:    []string{"--verbose", "false"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:    "required boolean not provided with default",
			expText: "hello true",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: true,
						Default: &Default{
							Value: true,
						},
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:     "required boolean not provided no default",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: no value for boolean flag "verbose"`,
			args:     nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:    "optional boolean not provided no default",
			expText: "hello false",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:     "no repeat boolean provided twice",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: multiple values set for boolean flag "verbose"`,
			args:     []string{"--verbose", "true", "--verbose", "true"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: false,
						Repeats: false,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:     "repeat boolean provided twice",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: multiple values set for boolean flag "verbose"`,
			args:     []string{"--verbose", "true", "--verbose", "true"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: false,
						Repeats: true,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose") // must use GetBools
					output = fmt.Sprintf("hello %t", verbose)
					return Success
				},
			},
		},
		{
			name:    "use equal sign true",
			expText: "ok true",
			expCode: Success,
			args:    []string{"--verbose=true"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: false,
						Repeats: false,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = "ok " + strconv.FormatBool(verbose)
					return Success
				},
			},
		},
		{
			name:    "use equal sign false",
			expText: "ok false",
			expCode: Success,
			args:    []string{"--verbose=false"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Require: false,
						Repeats: false,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBool("verbose")
					output = "ok " + strconv.FormatBool(verbose)
					return Success
				},
			},
		},
	}

	for _, tc := range cases {
		output = ""                    // reset for each case
		failure = new(strings.Builder) // reset for each case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
				Output:    failure,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.expText, output)
			must.Eq(t, tc.expCode, result)
			must.Eq(t, tc.expPanic, failure.String())
		})
	}
}

func TestComponent_GetBooleans(t *testing.T) {
	t.Parallel()

	var output string
	var failure *strings.Builder

	cases := []testCase{
		{
			name:    "repeated boolens provided no default",
			expText: "hello [true false true]",
			expCode: Success,
			args:    []string{"--verbose", "true", "--verbose", "false", "--verbose", "true"},
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Repeats: true,
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBools("verbose")
					output = fmt.Sprintf("hello %v", verbose)
					return Success
				},
			},
		},
		{
			name:    "repeated booleans not provided no default not required",
			expText: "hello []",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Repeats: true,
						Require: false,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBools("verbose")
					output = fmt.Sprintf("hello %v", verbose)
					return Success
				},
			},
		},
		{
			name:     "repeated booleans not provided no default required",
			expText:  "",
			expCode:  Failure,
			expPanic: `babycli: no value for boolean flag "verbose"`,
			args:     nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Repeats: true,
						Require: true,
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBools("verbose")
					output = fmt.Sprintf("hello %v", verbose)
					return Success
				},
			},
		},
		{
			name:    "repeated booleans not provided with default required",
			expText: "hello [true]",
			expCode: Success,
			args:    nil,
			root: &Component{
				Flags: Flags{
					{
						Type:    BooleanFlag,
						Long:    "verbose",
						Repeats: true,
						Require: true,
						Default: &Default{
							Value: true,
						},
					},
				},
				Function: func(c *Component) Code {
					verbose := c.GetBools("verbose")
					output = fmt.Sprintf("hello %v", verbose)
					return Success
				},
			},
		},
	}

	for _, tc := range cases {
		output = ""                    // reset for each case
		failure = new(strings.Builder) // reset for each case

		t.Run(tc.name, func(t *testing.T) {
			config := &Configuration{
				Arguments: tc.args,
				Top:       tc.root,
				Output:    failure,
			}
			c := New(config)
			result := c.Run()
			must.Eq(t, tc.expText, output)
			must.Eq(t, tc.expCode, result)
			must.Eq(t, tc.expPanic, failure.String())
		})
	}
}

func TestComponent_maybeSplit(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		arg  string
		exp  string
		push []string
	}{
		{
			name: "plain",
			arg:  "-name",
			exp:  "-name",
			push: nil,
		},
		{
			name: "split",
			arg:  "-name=bob",
			exp:  "-name",
			push: []string{"bob"},
		},
		{
			name: "quote",
			arg:  "'a=b'",
			exp:  "'a=b'",
			push: nil,
		},
		{
			name: "quote split",
			arg:  "-name='bob dylan'",
			exp:  "-name",
			push: []string{"'bob dylan'"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Component{
				args: stacks.Simple[string](),
			}
			result := c.maybeSplit(tc.arg)
			must.Eq(t, tc.exp, result)
			must.Eq(t, c.args.Size(), len(tc.push))
		})
	}
}
