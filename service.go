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
	`time`
	`golang.org/x/sys/windows/svc`
	`golang.org/x/sys/windows/svc/debug`
	`golang.org/x/sys/windows/svc/eventlog`
)

// https://golang.org/pkg/net/http/#Server.Shutdown
// https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2
// https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve

var elog debug.Log

type myservice struct{}

func (m *myservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {

	//const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}
	//fasttick := time.Tick(500 * time.Millisecond)
	//slowtick := time.Tick(2 * time.Second)
	//tick := fasttick
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:

	for {
		go func() {
			if err := conf.Server.ListenAndServe(); err != nil {
				elog.Error(1, err.Error())
			}
		}()

		//case <-tick:
		//	elog.Error(1, conf.Server.ListenAndServe().Error())
		//select {

		//case c := <-r:
		c := <-r

		switch c.Cmd {

		case svc.Interrogate:
			changes <- c.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			changes <- c.CurrentStatus

		// timeout could be given instead of nil as a https://golang.org/pkg/context/
		case svc.Stop, svc.Shutdown:
			if err := conf.Server.Shutdown(nil); err != nil {
				elog.Error(1, err.Error())
			}
			break loop

			//case svc.Pause:
				//changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				//tick = slowtick

			//case svc.Continue:
				//changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				//tick = fasttick

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
