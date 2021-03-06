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
	`time`

	`golang.org/x/sys/windows/svc`
	`golang.org/x/sys/windows/svc/debug`
	`golang.org/x/sys/windows/svc/eventlog`
)

var elog debug.Log

type myservice struct{}

func (this *myservice) Execute(args []string, req <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	Loop:

	for {
		go func() {
			if err := conf.Server.ListenAndServe(); err != nil {
				elog.Error(1, err.Error())
			}
		}()

		c := <-req

		switch c.Cmd {

		case svc.Interrogate:

			changes <- c.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			changes <- c.CurrentStatus

		case svc.Stop, svc.Shutdown:

			if err := conf.Server.Shutdown(nil); err != nil {
				elog.Error(1, err.Error())
			}

			break Loop

		default:
			elog.Error(1, fmt.Sprintf(`unexpected control request #%d`, c))
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {

	var err error

	if isDebug {
		elog = debug.New(name)
	} else if elog, err = eventlog.Open(name); err != nil {
		return
	}

	defer elog.Close()
	elog.Info(1, fmt.Sprintf(`starting %s service`, name))

	run := svc.Run

	if isDebug {
		run = debug.Run
	}

	if err = run(name, &myservice{}); err != nil {
		elog.Error(1, fmt.Sprintf(`%s service failed: %v`, name, err))
		return
	}

	elog.Info(1, fmt.Sprintf(`%s service stopped`, name))
}
