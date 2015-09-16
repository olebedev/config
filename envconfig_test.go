// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"testing"
)

var envYamlString = `
string: "val"
string2: "val2"
int: 1
int2: 2
float: 1.1
float2: 2.2
list: [1]
list2: [2]
map:
  key1: 1
map2:
  key1: 21
dev:
  string: "val11"
  int: 11
  float: 11.1
  list: [11]
  map:
    key1: 11
`

var envConfigTests = []struct {
	path string
	kind string
	want interface{}
	ok   bool
}{
	// Get from env.
	{"string", "String", "val11", true},
	{"int", "Int", 11, true},
	{"float", "Float64", 11.1, true},
	{"list", "List", []interface{}{11}, true},
	{"map", "Map", map[string]interface{}{"key1": 11}, true},

	// Get from root.
	{"string2", "String", "val2", true},
	{"int2", "Int", 2, true},
	{"float2", "Float64", 2.2, true},
	{"list2", "List", []interface{}{2}, true},
	{"map2", "Map", map[string]interface{}{"key1": 21}, true},
}

func TestEnvConfig(t *testing.T) {
	cfg, err := ParseYaml(envYamlString)
	if err != nil {
		t.Fatal(err)
	}
	envCfg := &EnvConfig{Env: "dev", Config: cfg}
	testEnvConfig(t, envCfg)
}

func testEnvConfig(t *testing.T, cfg *EnvConfig) {
Loop:
	for _, test := range envConfigTests {
		var got interface{}
		var err error
		switch test.kind {
		case "Bool":
			got, err = cfg.Bool(test.path)
		case "Float64":
			got, err = cfg.Float64(test.path)
		case "Int":
			got, err = cfg.Int(test.path)
		case "List":
			got, err = cfg.List(test.path)
		case "Map":
			got, err = cfg.Map(test.path)
		case "String":
			got, err = cfg.String(test.path)
		default:
			t.Errorf("Unsupported kind %q", test.kind)
			continue Loop
		}
		if test.ok {
			if err != nil {
				t.Errorf(`%s(%q) = "%v", got error: %v`, test.kind, test.path, test.want, err)
			} else {
				ok := false
				switch test.kind {
				case "List":
					ok = equalList(got, test.want)
				case "Map":
					ok = equalMap(got, test.want)
				default:
					ok = got == test.want
				}
				if !ok {
					t.Errorf(`%s(%q) = "%v", want "%v"`, test.kind, test.path, test.want, got)
				}
			}
		} else {
			if err == nil {
				t.Errorf("%s(%q): expected error", test.kind, test.path)
			}
		}
	}
}
