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

import `sort`

// -------
// Device.
// -------

// Device is shorthand for map of interface types indexed by string.
type Device map[string]interface{}

// VendorID returns the device vendor ID string.
func (this Device) VendorID() string {
	return this[`vendor_id`].(string)
}

// ProductID returns the device product ID string.
func (this Device) ProductID() string {
	return this[`product_id`].(string)
}

// SerialNumber returns the device serial number string.
func (this Device) SerialNumber() string {
	return this[`serial_number`].(string)
}

// --------
// Devices.
// --------

// Devices is a slice of Device objects.
type Devices []Device

// Sort sorts a collection devices by Vendor ID, Product ID, and Serial Num.
func (this *Devices) Sort() {
	sort.Sort(byDevice(*this))
}

// Add appends a new device to a collection of devices.
func (this *Devices) Add(dev Device) {
	*this = append(*this, dev)
}

// ---------
// byDevice.
// ---------

// byDevice allows custom sorting of a slice of devices.
type byDevice Devices

// Len returns the length of the slice.
func (this byDevice) Len() int {
	return len(this)
}

// Swap exchanges the indexes of two elements.
func (this byDevice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// Less provides the rules for sorting devices.
func (this byDevice) Less(i, j int) bool {

	vi, vj := this[i][`vendor_id`].(string), this[j][`vendor_id`].(string)
	pi, pj := this[i][`product_id`].(string), this[j][`product_id`].(string)
	si, sj := this[i][`serial_number`].(string), this[j][`serial_number`].(string)

	if vi != vj {
		return vi < vj
	} else if pi != pj {
		return pi < pj
	} else {
		return si < sj
	}
}
