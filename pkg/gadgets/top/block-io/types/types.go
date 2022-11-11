// Copyright 2019-2022 The Inspektor Gadget authors
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

package types

import (
	eventtypes "github.com/inspektor-gadget/inspektor-gadget/pkg/types"
)

type SortBy int

const (
	ALL SortBy = iota
	IO
	BYTES
	TIME
)

const (
	MaxRowsDefault  = 20
	IntervalDefault = 1
)

var SortByDefault = []string{"-io", "-bytes", "-us"}

const (
	IntervalParam = "interval"
	MaxRowsParam  = "max_rows"
	SortByParam   = "sort_by"
)

// Stats represents the operations performed on a single file
type Stats struct {
	eventtypes.CommonData

	Write      bool   `json:"write,omitempty" column:"write"`
	Major      int    `json:"major,omitempty" column:"major"`
	Minor      int    `json:"minor,omitempty" column:"minor"`
	Bytes      uint64 `json:"bytes,omitempty" column:"bytes"`
	MicroSecs  uint64 `json:"us,omitempty" column:"us"`
	Operations uint32 `json:"io,omitempty" column:"io"`
	MountNsID  uint64 `json:"mountnsid,omitempty" column:"mountnsid"`
	Pid        int32  `json:"pid,omitempty" column:"pid"`
	Comm       string `json:"comm,omitempty" column:"comm"`
}
