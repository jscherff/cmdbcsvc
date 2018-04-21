// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	`fmt`
	`net/http`
	`log`
	`html/template`
	`os`
	`path/filepath`
)

const configFile = `config.json`

var conf = new(Config)

// Config contains the configuration settings for the server.
type Config struct {
	Hostname	string
	Service		*Service
	Server		*http.Server
	Include		*Include
	Template	*template.Template
	Templates	[]string
	Resources	[]string
}

// Service contains the name and description of the windows service.
type Service struct {
	Name		string
	Description	string
}

// Include determines which devices are returned in the inventory.
type Include struct {
	VendorID	map[string]bool
	ProductID	map[string]map[string]bool
	Default		bool
}

// init loads the server configuration.
func init() {

	log.SetFlags(log.Flags() | log.Lshortfile)
	var appDir string

	// Determine the absolute path of the application directory.

	if path, err := filepath.Abs(os.Args[0]); err != nil {
		log.Fatal(err)
	} else {
		appDir = filepath.Dir(path)
	}

	// Prepend application directory to config file and load config.

	if err := load(conf, filepath.Join(appDir, configFile)); err != nil {
		log.Fatal(err)
	}

	// Obtain the hostname for templates.

	if hn, err := os.Hostname(); err != nil {
		log.Fatal(err)
	} else {
		conf.Hostname = hn
	}

	// Prepend application directory to template files.

	for index, file := range conf.Templates {
		conf.Templates[index] = filepath.Join(appDir, file)
	}

	// Load template files into template.

	if tmpl, err := template.ParseFiles(conf.Templates...); err != nil {
		log.Fatal(err)
	} else {
		conf.Template = tmpl
	}

	// Process resource paths and configure FileServer handler for each.

	for _, dir := range conf.Resources {

		osPath := filepath.Join(appDir, dir)
		urlPath := fmt.Sprintf(`/%s/`, dir)

		fs := http.FileServer(http.Dir(osPath))
		http.Handle(urlPath, http.StripPrefix(urlPath, fs))
	}
}
