// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"os"
	"reflect"
	"testing"
)

var yamlString = `
map:
  key0: true
  key1: false
  key2: "true"
  key3: "false"
  key4: 4.2
  key5: "4.2"
  key6: 42
  key7: "42"
  key8: value8
list:
  - true
  - false
  - "true"
  - "false"
  - 4.3
  - "4.3"
  - 43
  - "43"
  - item8
config:
  server:
    - www.google.com
    - www.cnn.com
    - www.example.com
  admin:
    - username: calvin
      password: yukon
    - username: hobbes
      password: tuna
messages:
  - |
    Welcome

    back!
  - >
    Farewell,

    my friend!
`

var configTests = []struct {
	path string
	kind string
	want interface{}
	ok   bool
}{
	// ok
	{"map.key0", "Bool", true, true},
	{"map.key0", "String", "true", true},
	// bad
	{"map.key0.foo", "Bool", "", false},
	{"map.key0", "Float64", "", false},
	{"map.key0", "Int", "", false},
	// ok
	{"map.key1", "Bool", false, true},
	{"map.key1", "String", "false", true},
	// bad
	{"map.key1", "Float64", "", false},
	{"map.key1", "Int", "", false},
	// ok
	{"map.key2", "Bool", true, true},
	{"map.key2", "String", "true", true},
	// bad
	{"map.key2", "Float64", "", false},
	{"map.key2", "Int", "", false},
	// ok
	{"map.key3", "Bool", false, true},
	{"map.key3", "String", "false", true},
	// bad
	{"map.key3", "Float64", "", false},
	{"map.key3", "Int", "", false},
	// ok
	{"map.key4", "Float64", 4.2, true},
	{"map.key4", "String", "4.2", true},
	// bad
	{"map.key4", "Bool", "", false},
	{"map.key4", "Int", "", false},
	// ok
	{"map.key5", "Float64", 4.2, true},
	{"map.key5", "String", "4.2", true},
	// bad
	{"map.key5", "Bool", "", false},
	{"map.key5", "Int", "", false},
	// ok
	{"map.key6", "Float64", float64(42), true},
	{"map.key6", "Int", 42, true},
	{"map.key6", "String", "42", true},
	// bad
	{"map.key6", "Bool", "", false},
	// ok
	{"map.key7", "Float64", float64(42), true},
	{"map.key7", "Int", 42, true},
	{"map.key7", "String", "42", true},
	// bad
	{"map.key7", "Bool", "", false},
	// ok
	{"map.key8", "String", "value8", true},
	// bad
	{"map.key8", "Bool", "", false},
	{"map.key8", "Float64", "", false},
	{"map.key8", "Int", "", false},
	// bad
	{"map.key9", "Bool", "", false},
	{"map.key9", "Float64", "", false},
	{"map.key9", "Int", "", false},
	{"map.key9", "String", "", false},

	// ok
	{"list.0", "Bool", true, true},
	{"list.0", "String", "true", true},
	// bad
	{"list.0", "Float64", "", false},
	{"list.0", "Int", "", false},
	// ok
	{"list.1", "Bool", false, true},
	{"list.1", "String", "false", true},
	// bad
	{"list.1", "Float64", "", false},
	{"list.1", "Int", "", false},
	// ok
	{"list.2", "Bool", true, true},
	{"list.2", "String", "true", true},
	// bad
	{"list.2", "Float64", "", false},
	{"list.2", "Int", "", false},
	// ok
	{"list.3", "Bool", false, true},
	{"list.3", "String", "false", true},
	// bad
	{"list.3", "Float64", "", false},
	{"list.3", "Int", "", false},
	// ok
	{"list.4", "Float64", 4.3, true},
	{"list.4", "String", "4.3", true},
	// bad
	{"list.4", "Bool", "", false},
	{"list.4", "Int", "", false},
	// ok
	{"list.5", "Float64", 4.3, true},
	{"list.5", "String", "4.3", true},
	// bad
	{"list.5", "Bool", "", false},
	{"list.5", "Int", "", false},
	// ok
	{"list.6", "Float64", float64(43), true},
	{"list.6", "Int", 43, true},
	{"list.6", "String", "43", true},
	// bad
	{"list.6", "Bool", "", false},
	// ok
	{"list.7", "Float64", float64(43), true},
	{"list.7", "Int", 43, true},
	{"list.7", "String", "43", true},
	// bad
	{"list.7", "Bool", "", false},
	// ok
	{"list.8", "String", "item8", true},
	// bad
	{"list.8", "Bool", "", false},
	{"list.8", "Float64", "", false},
	{"list.8", "Int", "", false},
	// bad
	{"list.9", "Bool", "", false},
	{"list.9", "Float64", "", false},
	{"list.9", "Int", "", false},
	{"list.9", "String", "", false},

	// ok
	{"config.server.0", "String", "www.google.com", true},
	{"config.server.1", "String", "www.cnn.com", true},
	{"config.server.2", "String", "www.example.com", true},
	// bad
	{"config.server.3", "Bool", "", false},
	{"config.server.3", "Float64", "", false},
	{"config.server.3", "Int", "", false},
	{"config.server.3", "String", "", false},

	// ok
	{"config.admin.0.username", "String", "calvin", true},
	{"config.admin.0.password", "String", "yukon", true},
	{"config.admin.1.username", "String", "hobbes", true},
	{"config.admin.1.password", "String", "tuna", true},
	// bad
	{"config.admin.0.country", "Bool", "", false},
	{"config.admin.0.country", "Float64", "", false},
	{"config.admin.0.country", "Int", "", false},
	{"config.admin.0.country", "String", "", false},

	// ok
	{"messages.0", "String", "Welcome\n\nback!\n", true},
	{"messages.1", "String", "Farewell,\nmy friend!\n", true},
	// bad
	{"messages.2", "Bool", "", false},
	{"messages.2", "Float64", "", false},
	{"messages.2", "Int", "", false},
	{"messages.2", "String", "", false},

	// ok
	{"config.server", "List", []interface{}{"www.google.com", "www.cnn.com", "www.example.com"}, true},
	{"config.admin.0", "Map", map[string]interface{}{"username": "calvin", "password": "yukon"}, true},
	{"config.admin.1", "Map", map[string]interface{}{"username": "hobbes", "password": "tuna"}, true},
}

func TestYamlConfig(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}
	str, err := RenderYaml(cfg.Root)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err = ParseYaml(str)
	if err != nil {
		t.Fatal(err)
	}
	testConfig(t, cfg)
}

func TestJsonConfig(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}
	str, err := RenderJson(cfg.Root)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err = ParseJson(str)
	if err != nil {
		t.Fatal(err)
	}
	testConfig(t, cfg)
}

func TestSet(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}
	val := "test"
	err = cfg.Set("map.key8", val)
	if v, _ := cfg.String("map.key8"); v != val {
		t.Errorf(`%s(%T) != "%s(%T)"`, v, v, val, val)
	}
}

func TestSetUnexistingValue(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}

	val := "test"

	err = cfg.Set("some.one", val)
	v, _ := cfg.String("some.one")
	expect(t, v, val)

	err = cfg.Set("some.thing.10", val)
	v, _ = cfg.String("some.thing.10")
	expect(t, v, val)
	// try to set by string key into slice
	expect(t, cfg.Set("some.thing.more", val) != nil, true)
}

func TestEnv(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}
	val := "test"
	cfg.Set("map.key8", "should be overwritten")
	os.Setenv("MAP_KEY8", val)
	cfg.Env()
	test, _ := cfg.String("map.key8")
	if test != val {
		t.Errorf(`"%s" != "%s"`, test, val)
	}
}

func TestEnvPrefix(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}
	val := "test"
	cfg.Set("map.key8", "should be overwritten")
	os.Setenv("PREFIX_MAP_KEY8", val)
	cfg.EnvPrefix("prefix")
	test, _ := cfg.String("map.key8")
	if test != val {
		t.Errorf(`"%s" != "%s"`, test, val)
	}
}

func TestFlag(t *testing.T) {
	cfg, err := ParseYaml(`
map:
  - listmap1:
      nested1: value1
      nested2: value2
    listmap2: value3
`)
	if err != nil {
		t.Fatal(err)
	}
	os.Args = append(os.Args, "-map-0-listmap2", "other")
	cfg.Flag()
	test, _ := cfg.String("map.0.listmap2")
	if test != "other" {
		t.Errorf(`"%s" != "%s"`, test, "other")
	}
}

func TestUMethods(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}

	// UString
	expect(t, cfg.UString("map.key8"), "value8")
	expect(t, cfg.UString("map.key8", "default"), "value8")
	expect(t, cfg.UString("map.undefined", "default"), "default")
	expect(t, cfg.UString("map.undefined"), "")

	// UBool
	expect(t, cfg.UBool("map.key0"), true)
	expect(t, cfg.UBool("map.key0", false), true)
	expect(t, cfg.UBool("map.undefined", true), true)
	expect(t, cfg.UBool("map.undefined"), false)

	// UFloat64
	expect(t, cfg.UFloat64("map.key4"), float64(4.2))
	expect(t, cfg.UFloat64("map.key4", float64(1)), float64(4.2))
	expect(t, cfg.UFloat64("map.undefined", float64(0.99)), float64(0.99))
	expect(t, cfg.UFloat64("map.undefined"), float64(0))

	// UInt
	expect(t, cfg.UInt("map.key6"), 42)
	expect(t, cfg.UInt("map.key6", 37), 42)
	expect(t, cfg.UInt("map.undefined", 37), 37)
	expect(t, cfg.UInt("map.undefined"), 0)

}

func TestCopy(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}

	cfg2, err := cfg.Copy()
	expect(t, err, nil)
	cfg2.Set("map.key6", 43)

	yaml1, _ := RenderYaml(cfg.Root)
	yaml2, _ := RenderYaml(cfg2.Root)

	expect(t, yaml2 == yaml1, false)

	cfg3, err := cfg.Copy("config", "server")
	expect(t, err, nil)
	cfg4, err := cfg.Copy("config.server")
	expect(t, err, nil)

	expect(t, cfg3.UString("0"), "www.google.com")
	expect(t, cfg4.UString("0"), "www.google.com")

	yaml3, _ := RenderYaml(cfg3.Root)
	yaml4, _ := RenderYaml(cfg4.Root)
	expect(t, yaml3, yaml4)
}

func TestExtendError(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}
	cfg2, err := ParseYaml(`
list:
  key0: true
map:
  - true
`)
	var nilCfg *Config
	extended, err := cfg.Extend(cfg2)
	expect(t, extended, nilCfg)
	expect(t, err.Error(), "Invalid list index at \"list.key0\"")
}

func TestExtend(t *testing.T) {
	cfg, err := ParseYaml(yamlString)
	if err != nil {
		t.Fatal(err)
	}
	cfg2, err := ParseYaml(`
map:
  key0: extend
list:
  - extend
`)

	extended, err := cfg.Extend(cfg2)
	expect(t, err, nil)
	// immutable
	expect(t, cfg.UBool("map.key0"), true)
	expect(t, cfg.UBool("list.0"), true)

	expect(t, cfg2.UString("map.key8", "not found"), "not found")
	expect(t, cfg2.UInt("list.8", 7), 7)

	// result
	expect(t, extended.UString("map.key0"), "extend")
	expect(t, extended.UString("map.key8"), "value8")
	expect(t, extended.UString("list.0"), "extend")
	expect(t, extended.UString("list.8"), "item8")
}

func TestComplexYamlKeys(t *testing.T) {
	cfg, err := ParseYaml(`
root:
  field1: value1
  field.something.2: value2
  "field number 3":
    field4: value3
  field.something.4:
    field5: value5
    field.6: value6
`)
	expect(t, err, nil)

	// result
	expect(t, cfg.UString("root.field1"), "value1")
	expect(t, cfg.UString("root.[field.something.2]"), "value2")
	expect(t, cfg.UString("root.field number 3.field4"), "value3")
	expect(t, cfg.UString("root.[field.something.4].field5"), "value5")
	expect(t, cfg.UString("root.[field.something.4].[field.6]"), "value6")
}

func testConfig(t *testing.T, cfg *Config) {
Loop:
	for _, test := range configTests {
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

func equalList(l1, l2 interface{}) bool {
	v1, ok1 := l1.([]interface{})
	v2, ok2 := l2.([]interface{})
	if !ok1 || !ok2 {
		return false
	}
	if len(v1) != len(v2) {
		return false
	}
	for k, v := range v1 {
		if v2[k] != v {
			return false
		}
	}
	return true
}

func equalMap(m1, m2 interface{}) bool {
	v1, ok1 := m1.(map[string]interface{})
	v2, ok2 := m2.(map[string]interface{})
	if !ok1 || !ok2 {
		return false
	}
	if len(v1) != len(v2) {
		return false
	}
	for k, v := range v1 {
		if v2[k] != v {
			return false
		}
	}
	return true
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
