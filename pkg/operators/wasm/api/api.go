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
	"fmt"
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

type logLevel uint32

// same as logrus, but hardcoded to avoid importing it
const (
	PanicLevel logLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func Log(level logLevel, args ...any) {
	message := fmt.Sprint(args...)
	xlog(uint32(level), uint64(toStringPtr(message)))
	runtime.KeepAlive(message)
}

func Logf(level logLevel, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	xlog(uint32(level), uint64(toStringPtr(message)))
	runtime.KeepAlive(message)
}

func Panic(params ...any) {
	Log(PanicLevel, params...)
}

func Panicf(fmt string, params ...any) {
	Logf(PanicLevel, fmt, params...)
}

func Fatal(params ...any) {
	Log(FatalLevel, params...)
}

func Fatalf(fmt string, params ...any) {
	Logf(FatalLevel, fmt, params...)
}

func Error(params ...any) {
	Log(ErrorLevel, params...)
}

func Errorf(fmt string, params ...any) {
	Logf(ErrorLevel, fmt, params...)
}

func Warn(params ...any) {
	Log(WarnLevel, params...)
}

func Warnf(fmt string, params ...any) {
	Logf(WarnLevel, fmt, params...)
}

func Info(params ...any) {
	Log(InfoLevel, params...)
}

func Infof(fmt string, params ...any) {
	Logf(InfoLevel, fmt, params...)
}

func Debug(params ...any) {
	Log(DebugLevel, params...)
}

func Debugf(fmt string, params ...any) {
	Logf(DebugLevel, fmt, params...)
}

func Trace(params ...any) {
	Log(TraceLevel, params...)
}

func Tracef(fmt string, params ...any) {
	Logf(TraceLevel, fmt, params...)
}

func GetDataSource(name string) DataSource {
	ret := getDataSource(uint64(toStringPtr(name)))
	runtime.KeepAlive(name)
	return ret
}

func NewDataSource(name string) DataSource {
	ret := newDataSource(uint64(toStringPtr(name)))
	runtime.KeepAlive(name)
	return DataSource(ret)
}

func (ds DataSource) Subscribe(cb func(DataSource, Data), priority uint32) {
	dsSubscriptionCtr++
	dsSubcriptions[dsSubscriptionCtr] = cb
	dataSourceSubscribe(uint32(ds), priority, dsSubscriptionCtr)
}

func (ds DataSource) NewData() Data {
	return Data(dataSourceNewData(uint32(ds)))
}

func (ds DataSource) EmitAndRelease(data Data) {
	dataSourceEmitAndRelease(uint32(ds), uint32(data))
}

func (ds DataSource) Release(data Data) {
	dataSourceRelease(uint32(ds), uint32(data))
}

func (ds DataSource) GetField(name string) Field {
	ret := dataSourceGetField(uint32(ds), uint64(toStringPtr(name)))
	runtime.KeepAlive(name)
	return Field(ret)
}

func (ds DataSource) AddField(name string) Field {
	ret := dataSourceAddField(uint32(ds), uint64(toStringPtr(name)))
	runtime.KeepAlive(name)
	return Field(ret)
}

func (f Field) String(data Data) string {
	str := fieldAccessorGetString(uint32(f), uint32(data))
	return stringPtr(str).String()
}

func (f Field) SetString(data Data, str string) {
	fieldAccessorSetString(uint32(f), uint32(data), uint64(toStringPtr(str)))
	runtime.KeepAlive(str)
}
