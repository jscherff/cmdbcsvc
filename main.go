// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Example service program that beeps.
//
// The program demonstrates how to create Windows service and
// install / remove it on a computer. It also shows how to
// stop / start / pause / continue any service, and how to
// write to event log. It also shows how to use debug
// facilities available in debug package.
//
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/svc"
)

const usageMsg =
	"%s\n\nUsage: %s <command>\n\twhere <command> is one of\n\t" +
	"install, remove, debug, start, or stop.\n"

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr, usageMsg, errmsg, os.Args[0])
	os.Exit(2)
}

func main() {

	var err error

	if isInteractive, err := svc.IsAnInteractiveSession(); err != nil {
		log.Fatalf("failed to determine if session is interactive: %v", err)
	} else if !isInteractive {
		runService(conf.Service.Name, false)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])

	switch cmd {

	case "debug":
		runService(conf.Service.Name, true)
		return

	case "install":
		err = installService(conf.Service.Name, conf.Service.Description)

	case "remove":
		err = removeService(conf.Service.Name)

	case "start":
		err = startService(conf.Service.Name)

	case "stop":
		err = controlService(conf.Service.Name, svc.Stop, svc.Stopped)

	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}

	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, conf.Service.Name, err)
	}

	return
}
