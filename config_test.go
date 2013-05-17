// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
    "testing"
)

var yamlString = `
map:
  key0: true
  key1: false
  key2: 4.2
  key3: 42
  key4: value5
list:
  - true
  - false
  - 4.3
  - 43
  - item4
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

var configTests = []struct{
    path string
    kind string
    want interface{}
    ok   bool
}{
    // ok
    {"map.key0", "Bool", true, true},
    {"map.key1", "Bool", false, true},
    {"map.key2", "Float64", 4.2, true},
    {"map.key3", "Int", 42, true},
    {"map.key4", "String", "value5", true},
    // bad
    {"map.key5", "Bool", "", false},
    {"map.key5", "Float64", "", false},
    {"map.key5", "Int", "", false},
    {"map.key5", "String", "", false},
    // ok
    {"list.0", "Bool", true, true},
    {"list.1", "Bool", false, true},
    {"list.2", "Float64", 4.3, true},
    {"list.3", "Int", 43, true},
    {"list.4", "String", "item4", true},
    // bad
    {"list.key5", "Bool", "", false},
    {"list.key5", "Float64", "", false},
    {"list.key5", "Int", "", false},
    {"list.key5", "String", "", false},
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
    testConfig(t, cfg, false)
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
    testConfig(t, cfg, true)
}

func testConfig(t *testing.T, cfg *Config, json bool) {
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
