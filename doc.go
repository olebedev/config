// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package config provides convenient access methods to configuration stored as
JSON or YAML.

Let's start with a simple YAML example:

	development:
	  database:
		host: localhost
	  users:
		- name: calvin
		  password: yukon
		- name: hobbes
		  password: tuna
	production:
	  database:
		host: 192.168.1.1

We can parse it using ParseYaml(), which will return a *Config instance on
success:

	cfg, err := ParseYaml(yamlString)
	// ...

An equivalent JSON configuration could be built using ParseJson():

	cfg, err := ParseJson(jsonString)
	// ...

From now, we can retrieve configuration values using a path in dotted notation.
So this:

	host, err := cfg.String("development.database.host")
	// ...

...will return the "localhost" value as stored in the configuration.

Besides String(), other types can be fetched directly: Bool(), Float64(),
Int(), Map() and List(). All these methods will return an error if the type
stored in the configuration doesn't match or can't be converted to the
requested type.

A nested configuration can be fetched using Get(). To get a new *Config
instance with a subset of the configuration, we can do:

	cfg, err := cfg.Get("development")
	// ...

Then the inner values are fetched relatively to the subset:

	host, err := cfg.String("database.host")
	// ...

For lists, the dotted path must use an index to refer to a specific value.
To retrieve the information from a user stored in the configuration above:

	user1, err := cfg.Map("development.users.0")
	// ...
	user2, err := cfg.Map("development.users.1")
	// ...

	// or...

	name1, err := cfg.String("development.users.0.name")
	// ...
	name2, err := cfg.String("development.users.1.name")
	// ...

We can render JSON or YAML calling the appropriate Render*() functions.
To render a configuration like the one used in these examples:

	cfg := map[string]interface{}{
		"development": map[string]interface{}{
			"database": map[string]interface{}{
				"host": "localhost",
			},
			"users": []interface{}{
				map[string]interface{}{
					"name":     "calvin",
					"password": "yukon",
				},
				map[string]interface{}{
					"name":     "hobbes",
					"password": "tuna",
				},
			},
		},
		"production": map[string]interface{}{
			"database": map[string]interface{}{
				"host": "192.168.1.1",
			},
		},
	}

	json, err := RenderJson(cfg)
	// ...

	// or...

	yaml, err := RenderYaml(cfg)
	// ...

This results in a configuration string to be stored in a file or database.
*/
package config
