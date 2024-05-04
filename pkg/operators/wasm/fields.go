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
	"github.com/inspektor-gadget/inspektor-gadget/pkg/gadget-service/api"
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
			wapi.GoModuleFunc(i.fieldAccessorGet),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI32, wapi.ValueTypeI32}, // Accessor, Data, Kind
			[]wapi.ValueType{wapi.ValueTypeI64},                                       // Val
		).
		Export("fieldAccessorGet")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.fieldAccessorSet),
			[]wapi.ValueType{
				wapi.ValueTypeI32, // Accessor
				wapi.ValueTypeI32, // Data
				wapi.ValueTypeI32, // Kind
				wapi.ValueTypeI64, // Value
			},
			[]wapi.ValueType{},
		).
		Export("fieldAccessorSet")
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
// - stack[2]: Kind
// Return value:
// - Uint8 value
func (i *wasmOperatorInstance) fieldAccessorGet(ctx context.Context, m wapi.Module, stack []uint64) {
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

	switch api.Kind(stack[2]) {
	case api.Kind_Int8:
		stack[0] = uint64(acc.Int8(data))
	case api.Kind_Int16:
		stack[0] = uint64(acc.Int16(data))
	case api.Kind_Int32:
		stack[0] = uint64(acc.Int32(data))
	case api.Kind_Int64:
		stack[0] = uint64(acc.Int64(data))
	case api.Kind_Uint8:
		stack[0] = uint64(acc.Uint8(data))
	case api.Kind_Uint16:
		stack[0] = uint64(acc.Uint16(data))
	case api.Kind_Uint32:
		stack[0] = uint64(acc.Uint32(data))
	case api.Kind_Uint64:
		stack[0] = uint64(acc.Uint64(data))
	case api.Kind_Float32:
		stack[0] = uint64(acc.Float32(data))
	case api.Kind_Float64:
		stack[0] = uint64(acc.Float64(data))
	//case api.Kind_Bool:
	//	stack[0] = uint64(acc.Bool(data))
	default:
		i.logger.Warnf("unknown field kind: %d", stack[2])
		stack[0] = 0
	}
}

// fieldAccessorSet saves a uint8 on the field
// Params:
// - stack[0]: Field handle
// - stack[1]: Data handle
// - stack[2]: Kind
// - stack[3]: Value to store
func (i *wasmOperatorInstance) fieldAccessorSet(ctx context.Context, m wapi.Module, stack []uint64) {
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

	switch api.Kind(stack[2]) {
	case api.Kind_Int8:
		acc.PutInt8(data, int8(stack[3]))
		stack[0] = uint64(acc.Int8(data))
	case api.Kind_Int16:
		acc.PutInt16(data, int16(stack[3]))
	case api.Kind_Int32:
		acc.PutInt32(data, int32(stack[3]))
	case api.Kind_Int64:
		acc.PutInt64(data, int64(stack[3]))
	case api.Kind_Uint8:
		acc.PutUint8(data, uint8(stack[3]))
	case api.Kind_Uint16:
		acc.PutUint16(data, uint16(stack[3]))
	case api.Kind_Uint32:
		acc.PutUint32(data, uint32(stack[3]))
	case api.Kind_Uint64:
		acc.PutUint64(data, uint64(stack[3]))
	// TODO: some missing types
	//case api.Kind_Float32:
	//	acc.PutFloat32(data, float32(stack[3]))
	//case api.Kind_Float64:
	//	acc.PutFloat64(data, float64(stack[3]))
	//case api.Kind_Bool:
	//	stack[0] = uint64(acc.Bool(data))
	default:
		i.logger.Warnf("unknown field kind: %d", stack[2])
		stack[0] = 0
	}
}
