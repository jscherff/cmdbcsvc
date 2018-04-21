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
	`golang.org/x/sys/windows/svc/mgr`
)

func startService(name string) (error) {

	m, err := mgr.Connect()

	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(name)

	if err != nil {
		return fmt.Errorf(`could not access service: %v`, err)
	}

	defer s.Close()

	err = s.Start(`is`, `manual-started`)

	if err != nil {
		return fmt.Errorf(`could not start service: %v`, err)
	}

	return nil
}

func controlService(name string, cmd svc.Cmd, to svc.State) (error) {

	m, err := mgr.Connect()

	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(name)

	if err != nil {
		return fmt.Errorf(`could not access service: %v`, err)
	}

	defer s.Close()

	status, err := s.Control(cmd)

	if err != nil {
		return fmt.Errorf(`could not send control=%d: %v`, cmd, err)
	}

	timeout := time.Now().Add(10 * time.Second)

	for status.State != to {

		if timeout.Before(time.Now()) {
			return fmt.Errorf(`timeout waiting for service to go to state=%d`, to)
		}

		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()

		if err != nil {
			return fmt.Errorf(`could not retrieve service status: %v`, err)
		}
	}

	return nil
}
