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
	"fmt"

	"github.com/inspektor-gadget/inspektor-gadget/pkg/operators/wasm/api"
)

//export init
func gadgetInit() int {
	api.Log("hello from wasm")
	ds := api.GetDataSource("dns")
	if ds == 0 {
		api.Log("failed to get datasource")
		return 1
	}

	nameF := ds.GetField("name")
	if nameF == 0 {
		api.Log("failed to get field")
		return 1
	}

	ds.Subscribe(func(source api.DataSource, data api.Data) {
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
				api.Log(fmt.Sprintf("invalid payload %+v\n", payload))
				return
			}
			i += length
		}

		nameF.SetString(data, str)
	}, 0)

	return 0
}

func main() {}
