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

func (i *wasmOperatorInstance) addDataSourceFuncs(env wazero.HostModuleBuilder) {
	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.newDataSource),
			[]wapi.ValueType{wapi.ValueTypeI64}, // DataSourceName
			[]wapi.ValueType{wapi.ValueTypeI32}, // DataSource
		).
		Export("newDataSource")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.getDataSource),
			[]wapi.ValueType{wapi.ValueTypeI64}, // DataSourceName
			[]wapi.ValueType{wapi.ValueTypeI32}, // DataSource
		).
		Export("getDataSource")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.dataSourceSubscribe),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI32, wapi.ValueTypeI64}, // DataSource, Priority, CallbackID
			[]wapi.ValueType{},
		).
		Export("dataSourceSubscribe")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.dataSourceGetField),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI64}, // DataSource, FieldName
			[]wapi.ValueType{wapi.ValueTypeI32},                    // Accessor
		).
		Export("dataSourceGetField")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.dataSourceAddField),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI64}, // DataSource, FieldName
			[]wapi.ValueType{wapi.ValueTypeI32},                    // Accessor
		).
		Export("dataSourceAddField")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.dataSourceNewData),
			[]wapi.ValueType{wapi.ValueTypeI32}, // DataSource
			[]wapi.ValueType{wapi.ValueTypeI32}, // Data
		).
		Export("dataSourceNewData")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.dataSourceEmitAndRelease),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI32}, // DataSource, Data
			[]wapi.ValueType{wapi.ValueTypeI32},                    // Error
		).
		Export("dataSourceEmitAndRelease")

	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.dataSourceRelease),
			[]wapi.ValueType{wapi.ValueTypeI32, wapi.ValueTypeI32}, // DataSource, Data
			[]wapi.ValueType{},
		).
		Export("dataSourceRelease")
}

// newDataSource creates a new datasource.
// Params:
// - stack[0] is the name of the datasource (string encoded)
// Return value:
// - datasource handle
func (i *wasmOperatorInstance) newDataSource(ctx context.Context, m wapi.Module, stack []uint64) {
	dsName, err := stringFromStack(m, stack[0])
	if err != nil {
		i.logger.Warnf("reading string from stack: %v", err)
		stack[0] = 0
		return
	}
	ds, err := i.gadgetCtx.RegisterDataSource(0, dsName)
	if err != nil {
		i.logger.Warnf("failed to register datasource: %v", err)
		stack[0] = 0
		return
	}
	stack[0] = wapi.EncodeU32(i.addHandle(ds))
}

// getDataSource returns a data source by its name.
// Params:
// - stack[0] is the name of the datasource (string encoded)
// Return value:
// - datasource handle
func (i *wasmOperatorInstance) getDataSource(ctx context.Context, m wapi.Module, stack []uint64) {
	dsName, err := stringFromStack(m, stack[0])
	if err != nil {
		i.logger.Warnf("reading string from stack: %v", err)
		stack[0] = 0
		return
	}
	ds := i.gadgetCtx.GetDataSources()[dsName]
	if ds == nil {
		i.logger.Debugf("datasource not found %q", dsName)
		stack[0] = 0
		return
	}
	stack[0] = wapi.EncodeU32(i.addHandle(ds))
}

// dataSourceGetField returns a handle to a data source.
// Params:
// - stack[0]: Datasource handle
// - stack[1]: Field name
// Return value:
// - Field handle
func (i *wasmOperatorInstance) dataSourceGetField(ctx context.Context, m wapi.Module, stack []uint64) {
	ds, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.DataSource)
	if !ok {
		i.logger.Warnf("datasource handle %d not found", stack[0])
		stack[0] = 0
		return
	}
	fieldName, err := stringFromStack(m, stack[1])
	if err != nil {
		i.logger.Warnf("reading string from stack: %v", err)
		stack[0] = 0
		return
	}
	acc := ds.GetField(fieldName)
	stack[0] = wapi.EncodeU32(i.addHandle(acc))
}

// dataSourceAddField add a field to the data source and returns its handle.
// Params:
// - stack[0]: Datasource handle
// - stack[1]: Field name
// Return value:
// - Field handle
// TODO: field type!
func (i *wasmOperatorInstance) dataSourceAddField(ctx context.Context, m wapi.Module, stack []uint64) {
	ds, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.DataSource)
	if !ok {
		i.logger.Warnf("datasource handle %d not found", stack[0])
		stack[0] = 0
		return
	}
	fieldName, err := stringFromStack(m, stack[1])
	if err != nil {
		i.logger.Warnf("reading string from stack: %v", err)
		stack[0] = 0
		return
	}
	acc, err := ds.AddField(fieldName)
	if err != nil {
		i.logger.Warnf("adding field %q to datasource %q: %v", fieldName, ds.Name(), err)
		stack[0] = 0
		return
	}
	stack[0] = wapi.EncodeU32(i.addHandle(acc))
}

// dataSourceSubscribe subscribes to the datasource.
// Params:
// - stack[0]: Datasource handle
// - stack[1]: Priority
// - stack[2]: Callback ID
// Return value:
// - 0 on success, 1 on error
func (i *wasmOperatorInstance) dataSourceSubscribe(ctx context.Context, m wapi.Module, stack []uint64) {
	ds, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.DataSource)
	if !ok {
		i.logger.Warnf("datasource handle %d not found", stack[0])
		stack[0] = 1
		return
	}
	prio := wapi.DecodeI32(stack[1])
	cbID := stack[2]
	cb := m.ExportedFunction("dsCallback")
	if cb == nil {
		i.logger.Warnf("module doesn't export dsCallback")
		stack[0] = 1
		return
	}
	ds.Subscribe(func(source datasource.DataSource, data datasource.Data) error {
		// TODO: can it be done outside? Why does DataFunc need the datasource again?
		tmpDS := i.addHandle(source)
		tmpData := i.addHandle(data)
		_, err := cb.Call(ctx, cbID, wapi.EncodeU32(tmpDS), wapi.EncodeU32(tmpData))
		i.delHandle(tmpDS)
		i.delHandle(tmpData)
		return err
	}, int(prio))

	stack[0] = 0
}

// dataSourceNewData allocates and returns a handle to the data.
// Params:
// - stack[0]: Datasource handle
// Return value:
// - data handle, 0 on error
func (i *wasmOperatorInstance) dataSourceNewData(ctx context.Context, m wapi.Module, stack []uint64) {
	ds, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.DataSource)
	if !ok {
		i.logger.Warnf("datasource handle %d not found", stack[0])
		stack[0] = 0
		return
	}
	data := ds.NewData()
	stack[0] = wapi.EncodeU32(i.addHandle(data))
}

// dataSourceEmitAndRelease emits and releases the data.
// Params:
// - stack[0]: Data handle
// Return value:
// - 0 on success, 1 on error
func (i *wasmOperatorInstance) dataSourceEmitAndRelease(ctx context.Context, m wapi.Module, stack []uint64) {
	ds, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.DataSource)
	if !ok {
		i.logger.Warnf("datasource handle %d not found", stack[0])
		stack[0] = 1
		return
	}
	data, ok := i.getHandle(wapi.DecodeU32(stack[1])).(datasource.Data)
	if !ok {
		i.gadgetCtx.Logger().Warnf("data handle %d not found", stack[1])
		stack[0] = 1
		return
	}
	ds.EmitAndRelease(data)
	stack[0] = 0
}

// dataSourceRelease releases the data.
// Params:
// - stack[0]: Data handle
// Return value:
// - 0 on success, 1 on error
func (i *wasmOperatorInstance) dataSourceRelease(ctx context.Context, m wapi.Module, stack []uint64) {
	ds, ok := i.getHandle(wapi.DecodeU32(stack[0])).(datasource.DataSource)
	if !ok {
		i.logger.Warnf("datasource handle %d not found", stack[0])
		stack[0] = 1
		return
	}
	data, ok := i.getHandle(wapi.DecodeU32(stack[1])).(datasource.Data)
	if !ok {
		i.gadgetCtx.Logger().Warnf("data handle %d not found", stack[1])
		stack[0] = 1
		return
	}
	ds.Release(data)
	stack[0] = 0
}
