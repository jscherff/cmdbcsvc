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

import `time`

type Inventory struct {
	Hostname string
	Devices Devices
	DateStamp string
}

// NewInventory scans for devices and generates a new inventory.
func NewInventory() (*Inventory, error) {

	if devices, err := scan(conf.Include); err != nil {
		return nil, err
	} else {
		devices.Sort()
		this := &Inventory{
			Hostname: conf.Hostname,
			Devices: devices,
			DateStamp: time.Now().Format(time.RFC3339),
		}
		return this, nil
	}
}
