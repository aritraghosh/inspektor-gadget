// Copyright 2024 The Inspektor Gadget authors
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
	"github.com/inspektor-gadget/inspektor-gadget/pkg/operators/wasm/api"
)

//export init
func gadgetInit() int {
	api.Info("hello from wasm")
	api.Warn("hello from wasm warn")

	ds := api.GetDataSource("dns")
	if ds == 0 {
		api.Warn("failed to get datasource")
		return 1
	}

	nameF := ds.GetField("name")
	if nameF == 0 {
		api.Warn("failed to get field")
		return 1
	}

	uidF := ds.GetField("uid")
	if nameF == 0 {
		api.Warn("failed to get field")
		return 1
	}

	ds.Subscribe(func(source api.DataSource, data api.Data) {
		uidF.SetUint32(data, 1234)

		payload := nameF.String(data)

		var str string
		for i := 0; i < len(payload); i++ {
			length := int(payload[i])
			if length == 0 {
				break
			}
			if i+1+length < len(payload) {
				str += string(payload[i+1:i+1+length]) + "."
			} else {
				api.Warn("invalid payload %+v", payload)
				return
			}
			i += length
		}

		nameF.SetString(data, str)
	}, 0)

	return 0
}

func main() {}
