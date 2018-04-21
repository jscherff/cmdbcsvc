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
)

func init() {
	http.HandleFunc(`/`, InventoryHandler)
}

const (
	inventoryTitle = `HID Inventory`
	inventoryCaption = `USB Human Input Device Inventory for %s`
)

// InventoryHandler responds to HTTP requests for device inventories.
func InventoryHandler(w http.ResponseWriter, r *http.Request) {

	caption := fmt.Sprintf(inventoryCaption, conf.Hostname)

	if inventory, err := NewInventory(inventoryTitle, caption); err != nil {
		panic(err)
	} else {
		w.Header().Set(`Cache-Control`, `no-store, must-revalidate`)
		w.Header().Set(`Expires`, `0`)
		conf.Template.Execute(w, inventory)
	}
}
