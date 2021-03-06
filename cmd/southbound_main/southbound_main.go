// Copyright 2019-present Open Networking Foundation.
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
	"fmt"
	"github.com/onosproject/onos-config/pkg/southbound"
	"github.com/onosproject/onos-config/pkg/southbound/topocache"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/client"
)

func main() {
	device := topocache.Device{
		Addr:     "localhost:10161",
		Target:   "Test-onos-config",
		// Loaded from default-certificates.go
		//CaPath:   "/Users/andrea/go/src/github.com/opennetworkinglab/onos-config/tools/test/devicesim/certs/onfca.crt",
		//CertPath: "/Users/andrea/go/src/github.com/opennetworkinglab/onos-config/tools/test/devicesim/certs/client1.crt",
		//KeyPath:  "/Users/andrea/go/src/github.com/opennetworkinglab/onos-config/tools/test/devicesim/certs/client1.key",
		Timeout:  time.Second * 10,
	}
	target, err := southbound.GetTarget(southbound.Key{Key: device.Addr})
	if err != nil {
		fmt.Println("Creating device for addr: ", device.Addr)
		target, _, err = southbound.ConnectTarget(device)
		if err != nil {
			fmt.Println("Error ", target, err)
		}
	}
	targetExists, err := southbound.GetTarget(southbound.Key{Key: device.Addr})
	if reflect.DeepEqual(target, targetExists) {
		fmt.Println("Target reusal works")
	}
	request := ""
	capResponse, capErr := southbound.CapabilitiesWithString(target, request)
	if capErr != nil {
		fmt.Println("Error ", target, err)
	}
	capResponseString := proto.MarshalTextString(capResponse)
	fmt.Println("Capabilities: ", capResponseString)
	request = "path: <elem: <name: 'system'> elem:<name:'config'> elem: <name: 'hostname'>>"
	getResponse, getErr := southbound.GetWithString(target, request)
	if getErr != nil {
		fmt.Println("Error ", target, err)
	}
	getResponseString := proto.MarshalTextString(getResponse)
	fmt.Println("get: ", getResponseString)

	//"Subscribe"
	options := &southbound.SubscribeOptions{
		UpdatesOnly:       false,
		Prefix:            "",
		Mode:              "stream",
		StreamMode:        "target_defined",
		SampleInterval:    15,
		HeartbeatInterval: 15,
		Paths:             nil,
		Origin:            "",
	}
	req, err := southbound.NewSubscribeRequest(options)

	handler := func(n client.Notification) error {
		switch v := n.(type) {
		case client.Update:
			fmt.Println(v.Path)
			//TODO
		case client.Delete:
			//TODO
		case client.Sync:
			//TODO
		case client.Error:
			//TODO
		}
		return nil
	}
	southbound.Subscribe(target, req, handler)
}
