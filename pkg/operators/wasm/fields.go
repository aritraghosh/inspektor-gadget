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

package wasm

import (
	"context"

	"github.com/tetratelabs/wazero"
	wapi "github.com/tetratelabs/wazero/api"

	"github.com/inspektor-gadget/inspektor-gadget/pkg/datasource"
)

func (i *wasmOperatorInstance) addFieldFuncs(env wazero.HostModuleBuilder) {
	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.fieldAccessorGetString),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI32}, // Accessor, Data
			[]wapi.ValueType{wapi.ValueTypeI64},                    // String
		).
		Export("fieldAccessorGetString")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.fieldAccessorSetString),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI32, wapi.ValueTypeI64}, // Accessor, Data
			[]wapi.ValueType{},
		).
		Export("fieldAccessorSetString")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.fieldAccessorGetUint8),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI32}, // Accessor, Data
			[]wapi.ValueType{wapi.ValueTypeI32},                    // Val
		).
		Export("fieldAccessorGetUint8")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.fieldAccessorSetUint8),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI32, wapi.ValueTypeI32}, // Accessor, Data
			[]wapi.ValueType{},
		).
		Export("fieldAccessorSetUint8")
}

// fieldAccessorGetString returns the field as a string
// Params:
// - stack[0]: Field handle
// - stack[1]: Data handle
// Return value:
// - String handle in success, 0 in case of error
func (i *wasmOperatorInstance) fieldAccessorGetString(ctx context.Context, m wapi.Module, stack []uint64) {
	acc, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.FieldAccessor)
	if !ok {
		i.logger.Warnf("field handle %d not found", stack[0])
		stack[0] = 0
		return
	}
	data, ok := i.getHandle(wapi.DecodeU32(stack[1])).(datasource.Data)
	if !ok {
		i.logger.Warnf("data handle %d not found", stack[1])
		stack[0] = 0
		return
	}

	str := []byte(acc.String(data))

	res, err := i.guestMalloc.Call(ctx, uint64(len(str)))
	if err != nil {
		i.logger.Warnf("malloc failed: %v", err)
		stack[0] = 0
		return
	}

	if !m.Memory().Write(uint32(res[0]), str) {
		i.logger.Warnf("out of memory write")
		stack[0] = 0
		return
	}

	stack[0] = uint64(len(str))<<32 | uint64(res[0])
}

// fieldAccessorSetString saves a string on the field
// Params:
// - stack[0]: Field handle
// - stack[1]: Data handle
// - stack[2]: String handle
// Return value:
// - 0 in success, 1 in case of error
func (i *wasmOperatorInstance) fieldAccessorSetString(ctx context.Context, m wapi.Module, stack []uint64) {
	acc, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.FieldAccessor)
	if !ok {
		i.logger.Warnf("field handle %d not found", stack[0])
		stack[0] = 0
		return
	}
	data, ok := i.getHandle(wapi.DecodeU32(stack[1])).(datasource.Data)
	if !ok {
		i.logger.Warnf("data handle %d not found", stack[1])
		stack[0] = 0
		return
	}

	str, err := stringFromStack(m, stack[2])
	if err != nil {
		i.logger.Warnf("reading string from stack: %v", err)
		stack[0] = 0
		return
	}

	// TODO: is this logic correct?
	buf := []byte(str)

	// fill the string with 0s if it's a static field
	s := acc.Size()
	if s != 0 {
		if len(buf) > int(s) {
			i.logger.Warnf("string too long: %d > %d", len(buf), s)
			stack[0] = 0
			return
		}

		buf = append(buf, make([]byte, int(s)-len(buf))...)
	}

	if err := acc.Set(data, buf); err != nil {
		i.logger.Warnf("setting string failed: %v", err)
		stack[0] = 0
		return
	}
}

// fieldAccessorGetUint8 returns the field as a uint8
// Params:
// - stack[0]: Field handle
// - stack[1]: Data handle
// Return value:
// - Uint8 value
func (i *wasmOperatorInstance) fieldAccessorGetUint8(ctx context.Context, m wapi.Module, stack []uint64) {
	acc, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.FieldAccessor)
	if !ok {
		i.logger.Warnf("field handle %d not found", stack[0])
		stack[0] = 0
		return
	}
	data, ok := i.getHandle(wapi.DecodeU32(stack[1])).(datasource.Data)
	if !ok {
		i.logger.Warnf("data handle %d not found", stack[1])
		stack[0] = 0
		return
	}

	stack[0] = uint64(acc.Uint8(data))
}

// fieldAccessorSetUint8 saves a uint8 on the field
// Params:
// - stack[0]: Field handle
// - stack[1]: Data handle
// - stack[2]: Value to store
func (i *wasmOperatorInstance) fieldAccessorSetUint8(ctx context.Context, m wapi.Module, stack []uint64) {
	acc, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.FieldAccessor)
	if !ok {
		i.logger.Warnf("field handle %d not found", stack[0])
		stack[0] = 0
		return
	}
	data, ok := i.getHandle(wapi.DecodeU32(stack[1])).(datasource.Data)
	if !ok {
		i.logger.Warnf("data handle %d not found", stack[1])
		stack[0] = 0
		return
	}

	acc.PutUint8(data, uint8(stack[2]))
}

// TODO: complete with other data types
