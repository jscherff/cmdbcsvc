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
	"Usage:\t%s <command>\n\t" +
	"Where <command> is one of:\n\t" +
	"install, remove, start, stop, debug, version, or help\n"

var (
	program string = filepath.Base(os.Args[0])
	version string = `undefined`
)

func showUsage() {
	fmt.Fprintf(os.Stderr, usageMsg, program)
}

func showVersion() {
	fmt.Fprintf(os.Stdout, `%s %s`, program, version)
}

func main() {

	if len(os.Args) == 2 {
		processCommand(os.Args[1])
		return
	}

	interactive, err := svc.IsAnInteractiveSession()

	if err != nil {
		log.Fatalf(`failed to determine if session is interactive: %v`, err)
	} else if interactive {
		log.Printf(`invalid command line for interactive session`)
		showUsage(); os.Exit(2)
	} else {
		runService(conf.Service.Name, false)
	}
}

func processCommand(cmd string) {

	var err error
	cmd = strings.ToLower(cmd)

	switch cmd {

	case `install`:
		err = installService(conf.Service.Name, conf.Service.Description)

	case `remove`:
		err = removeService(conf.Service.Name)

	case `start`:
		err = startService(conf.Service.Name)

	case `stop`:
		err = controlService(conf.Service.Name, svc.Stop, svc.Stopped)

	case `debug`:
		runService(conf.Service.Name, true)

	case `version`:
		showVersion()

	case `help`:
		showUsage()

	default:
		log.Printf(`invalid command %s`, cmd)
		showUsage(); os.Exit(2)
	}

	if err != nil {
		log.Fatalf(`failed to %s %s: %v`, cmd, conf.Service.Name, err)
	}
}
