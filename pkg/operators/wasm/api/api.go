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

package api

// TODO: is it possible to make it work without cgo?

// #include <stdlib.h>
import "C"

import (
	"runtime"
	"strings"
	"unsafe"
)

type stringPtr uint64

// toStringPtr gets the the pointer and length of the string as a uint64.
// Callers must use runtime.KeepAlive on the input string to ensure it is not
// garbage collected
func toStringPtr(s string) stringPtr {
	unsafePtr := unsafe.Pointer(unsafe.StringData(s))
	return stringPtr(uint64(len(s))<<32 | uint64(uintptr(unsafePtr)))
}

func (s stringPtr) String() string {
	if s == 0 {
		return ""
	}
	// create a string that users the pointer as storage
	orig := unsafe.String((*byte)(unsafe.Pointer(uintptr(s&0xFFFFFFFF))), int(s>>32))
	// clone it
	ret := strings.Clone(orig)
	// free the original pointer
	C.free(unsafe.Pointer(uintptr(s & 0xFFFFFFFF)))
	// return the cloned string
	return ret
}

var (
	dsSubscriptionCtr = uint64(0)
	dsSubcriptions    = map[uint64]func(DataSource, Data){}
)

//export dsCallback
func dsCallback(cbID uint64, ds uint32, data uint32) {
	cb, ok := dsSubcriptions[cbID]
	if !ok {
		return
	}
	cb(DataSource(ds), Data(data))
}

type (
	DataSource uint32
	Field      uint32
	Data       uint32
)

func Log(message string) {
	xlog(toStringPtr(message))
	runtime.KeepAlive(message)
}

func GetDataSource(name string) DataSource {
	ret := getDataSource(toStringPtr(name))
	runtime.KeepAlive(name)
	return ret
}

func NewDataSource(name string) DataSource {
	ret := newDataSource(toStringPtr(name))
	runtime.KeepAlive(name)
	return ret
}

func (ds DataSource) Subscribe(cb func(DataSource, Data), priority uint32) {
	dsSubscriptionCtr++
	dsSubcriptions[dsSubscriptionCtr] = cb
	dataSourceSubscribe(ds, priority, dsSubscriptionCtr)
}

func (ds DataSource) NewData() Data {
	return dataSourceNewData(ds)
}

func (ds DataSource) EmitAndRelease(data Data) {
	dataSourceEmitAndRelease(ds, data)
}

func (ds DataSource) Release(data Data) {
	dataSourceRelease(ds, data)
}

func (ds DataSource) GetField(name string) Field {
	ret := dataSourceGetField(ds, toStringPtr(name))
	runtime.KeepAlive(name)
	return ret
}

func (ds DataSource) AddField(name string) Field {
	ret := dataSourceAddField(ds, toStringPtr(name))
	runtime.KeepAlive(name)
	return ret
}

func (f Field) String(data Data) string {
	str := fieldAccessorGetString(f, data)
	return str.String()
}

func (f Field) SetString(data Data, str string) {
	fieldAccessorSetString(f, data, toStringPtr(str))
	runtime.KeepAlive(str)
}
