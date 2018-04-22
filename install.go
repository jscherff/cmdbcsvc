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
	`os`
	`path/filepath`

	`golang.org/x/sys/windows/svc/eventlog`
	`golang.org/x/sys/windows/svc/mgr`
)

func exePath() (string, error) {

	check := func(path string) error {

		if fileInfo, err := os.Stat(path); err != nil {
			return err
		} else if fileInfo.Mode().IsDir() {
			return fmt.Errorf(`%s is directory`, path)
		}

		return nil
	}

	prog := os.Args[0]
	path, err := filepath.Abs(prog)

	if err != nil {
		return ``, err
	}

	if err = check(path); err == nil {
		return path, nil
	}

	if filepath.Ext(path) == `` {
		path += `.exe`
	}

	if err = check(path); err == nil {
		return path, nil
	}

	return ``, err
}

func installService(name, desc string) (error) {

	path, err := exePath()

	if err != nil {
		return err
	}

	m, err := mgr.Connect()

	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(name)

	if err == nil {
		s.Close()
		return fmt.Errorf(`service %s already exists`, name)
	}

	conf := mgr.Config{StartType: mgr.StartAutomatic, DisplayName: desc}
	s, err = m.CreateService(name, path, conf, `is`, `auto-started`)

	if err != nil {
		return err
	}

	defer s.Close()

	err = eventlog.InstallAsEventCreate(
		name,
		eventlog.Error|eventlog.Warning|eventlog.Info,
	)

	if err != nil {
		s.Delete()
		return fmt.Errorf(`SetupEventLogSource() failed: %s`, err)
	}

	return nil
}

func removeService(name string) (error) {

	m, err := mgr.Connect()

	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(name)

	if err != nil {
		return fmt.Errorf(`service %s is not installed`, name)
	}

	defer s.Close()

	err = s.Delete()

	if err != nil {
		return err
	}

	err = eventlog.Remove(name)

	if err != nil {
		return fmt.Errorf(`RemoveEventLogSource() failed: %s`, err)
	}

	return nil
}
