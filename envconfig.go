// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import ()

// Config ---------------------------------------------------------------------

// EnvConfig provides convenient access to settings for an environment.
// It will first check "ENV.path", and if not found, return "PATH".
// This makes it easy to provide default settings which can be overwritten
// for each environment.
//
// To use it, initialize your config first.
// Then create a new EnvConfig like this:
// envConfig := config.EnvConfig{Env: "production", Config: &conf}
type EnvConfig struct {
	Env    string
	Config *Config
}

// Get returns a nested config according to a dotted path.
func (c *EnvConfig) Get(path string) (*Config, error) {
	return c.Config.Get(path)
}

// Set a nested config according to a dotted path.
func (c *EnvConfig) Set(path string, val interface{}) error {
	return c.Config.Set(path, val)
}

// Bool returns a bool according to a dotted path.
func (c *EnvConfig) Bool(path string) (bool, error) {
	val, err := c.Config.Bool(c.Env + "." + path)
	if err != nil {
		val, err = c.Config.Bool(path)
	}
	return val, err
}

// UBool retirns a bool according to a dotted path or default value or false.
func (c *EnvConfig) UBool(path string, defaults ...bool) bool {
	value, err := c.Bool(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return false
}

// Float64 returns a float64 according to a dotted path.
func (c *EnvConfig) Float64(path string) (float64, error) {
	val, err := c.Config.Float64(c.Env + "." + path)
	if err != nil {
		val, err = c.Config.Float64(path)
	}
	return val, err
}

// UFloat64 returns a float64 according to a dotted path or default value or 0.
func (c *EnvConfig) UFloat64(path string, defaults ...float64) float64 {
	value, err := c.Float64(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return float64(0)
}

// Int returns an int according to a dotted path.
func (c *EnvConfig) Int(path string) (int, error) {
	val, err := c.Config.Int(c.Env + "." + path)
	if err != nil {
		val, err = c.Config.Int(path)
	}
	return val, err
}

// UInt returns an int according to a dotted path or default value or 0.
func (c *EnvConfig) UInt(path string, defaults ...int) int {
	value, err := c.Int(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return 0
}

// List returns a []interface{} according to a dotted path.
func (c *EnvConfig) List(path string) ([]interface{}, error) {
	val, err := c.Config.List(c.Env + "." + path)
	if err != nil {
		val, err = c.Config.List(path)
	}
	return val, err
}

// UList returns a []interface{} according to a dotted path or defaults or []interface{}.
func (c *EnvConfig) UList(path string, defaults ...[]interface{}) []interface{} {
	value, err := c.List(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return make([]interface{}, 0)
}

// Map returns a map[string]interface{} according to a dotted path.
func (c *EnvConfig) Map(path string) (map[string]interface{}, error) {
	val, err := c.Config.Map(c.Env + "." + path)
	if err != nil {
		val, err = c.Config.Map(path)
	}
	return val, err
}

// UMap returns a map[string]interface{} according to a dotted path or default or map[string]interface{}.
func (c *EnvConfig) UMap(path string, defaults ...map[string]interface{}) map[string]interface{} {
	value, err := c.Map(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return map[string]interface{}{}
}

// String returns a string according to a dotted path.
func (c *EnvConfig) String(path string) (string, error) {
	val, err := c.Config.String(c.Env + "." + path)
	if err != nil {
		val, err = c.Config.String(path)
	}
	return val, err
}

// UString returns a string according to a dotted path or default or "".
func (c *EnvConfig) UString(path string, defaults ...string) string {
	value, err := c.String(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return ""
}
