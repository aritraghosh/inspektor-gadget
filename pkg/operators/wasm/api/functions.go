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

// This file contains all functions that are exported to the wasm guest module.
// Keep this aligned with pkg/operators/wasm/wasm.go

package api

//go:wasmimport env xlog
func xlog(s stringPtr)

//go:wasmimport env newDataSource
func newDataSource(name stringPtr) DataSource

//go:wasmimport env getDataSource
func getDataSource(name stringPtr) DataSource

//go:wasmimport env dataSourceSubscribe
func dataSourceSubscribe(ds DataSource, prio uint32, cb uint64)

//go:wasmimport env dataSourceGetField
func dataSourceGetField(ds DataSource, name stringPtr) Field

//go:wasmimport env dataSourceAddField
func dataSourceAddField(ds DataSource, name stringPtr) Field

//go:wasmimport env getField
func getField(ds DataSource) Field

//go:wasmimport env dataSourceNewData
func dataSourceNewData(ds DataSource) Data

//go:wasmimport env dataSourceEmitAndRelease
func dataSourceEmitAndRelease(ds DataSource, data Data)

//go:wasmimport env dataSourceRelease
func dataSourceRelease(ds DataSource, data Data)

//go:wasmimport env fieldAccessorGetString
func fieldAccessorGetString(acc Field, data Data) stringPtr

//go:wasmimport env fieldAccessorSetString
func fieldAccessorSetString(acc Field, data Data, str stringPtr)

////go:wasmimport env mfree
//func mfree(uint32)
//
////go:wasmimport env freeHost
//func freeHost(entry uint32)
//
