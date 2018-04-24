// Copyright 2017 John Scherff | Copyright 2012 The Go Authors
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
	`log`
	`os`
	`path/filepath`
	`strings`

	`golang.org/x/sys/windows/svc`
)

const usageMsg =
	"%s\n\nUsage: %s <command>\n\twhere <command> is one of\n\t" +
	"version, install, remove, start, stop, restart, or debug\n"

var (
	program string = filepath.Base(os.Args[0])
	version string = `undefined`
)

func showUsage(errMsg string) {
	fmt.Fprintf(os.Stderr, usageMsg, errMsg, os.Args[0])
	os.Exit(2)
}

func showVersion() {
	fmt.Fprintf(os.Stdout, `%s %s`, program, version)
	os.Exit(0)
}

func main() {

	switch len(os.Args) {

	case 1:
		handleStartup()

	case 2:
		processCommand(os.Args[1])

	default:
		showUsage(`invalid command line`)
	}
}

func processCommand(cmd string) {

	var err error
	cmd = strings.ToLower(cmd)

	switch cmd {

	case `version`:
		showVersion()

	case `install`:
		err = installService(conf.Service.Name, conf.Service.Description)

	case `remove`:
		err = removeService(conf.Service.Name)

	case `start`:
		err = startService(conf.Service.Name)

	case `stop`:
		err = controlService(conf.Service.Name, svc.Stop, svc.Stopped)

	case `restart`:
		err = controlService(conf.Service.Name, svc.Stop, svc.Stopped)
		if err == nil { err = startService(conf.Service.Name) }

	case `debug`:
		runService(conf.Service.Name, true)

	default:
		showUsage(fmt.Sprintf(`invalid command %s`, cmd))
	}

	if err != nil {
		log.Fatalf(`failed to %s %s: %v`, cmd, conf.Service.Name, err)
	}
}

func handleStartup() {

	if interactive, err := svc.IsAnInteractiveSession(); err != nil {
		log.Fatalf(`failed to determine if session is interactive: %v`, err)
	} else if !interactive {
		runService(conf.Service.Name, false)
	}
}
