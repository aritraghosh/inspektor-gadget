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
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	log "github.com/sirupsen/logrus"
	"github.com/tetratelabs/wazero"
	wapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"

	"github.com/inspektor-gadget/inspektor-gadget/pkg/gadget-service/api"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/logger"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/oci"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/operators"
)

const (
	ParamGlobalAllowHostFS = "wasm-global-allow-host-fs" // TODO
	ParamAllowHostFS       = "wasm-allow-host-fs"

	wasmObjectMediaType = "application/vnd.gadget.wasm.program.v1+binary"
)

type wasmOperator struct{}

func (w *wasmOperator) Name() string {
	return "wasm"
}

func (w *wasmOperator) Description() string {
	return "TODO"
}

func (w *wasmOperator) InstantiateImageOperator(
	gadgetCtx operators.GadgetContext,
	desc ocispec.Descriptor,
	paramValues api.ParamValues,
) (
	operators.ImageOperatorInstance, error,
) {
	return &wasmOperatorInstance{
		desc:        desc,
		gadgetCtx:   gadgetCtx,
		handleMap:   map[uint32]any{},
		allowHostFS: paramValues[ParamAllowHostFS] == "true",
		logger:      gadgetCtx.Logger(),
	}, nil
}

type wasmOperatorInstance struct {
	desc      ocispec.Descriptor
	rt        wazero.Runtime
	gadgetCtx operators.GadgetContext
	mod       wapi.Module
	// Golang objects are exposed to the wasm module by using a handleID
	handleMap  map[uint32]any
	handleCtr  uint32
	handleLock sync.RWMutex

	allowHostFS bool
	logger      logger.Logger

	// malloc function exported by the guess
	guestMalloc wapi.Function
}

func (w *wasmOperatorInstance) Name() string {
	return "wasm"
}

func (w *wasmOperatorInstance) Prepare(gadgetCtx operators.GadgetContext) error {
	err := w.init(gadgetCtx)
	if err != nil {
		return fmt.Errorf("initializing wasm: %w", err)
	}
	err = w.Init(gadgetCtx)
	if err != nil {
		return fmt.Errorf("initializing wasm guest: %w", err)
	}

	return nil
}

func (w *wasmOperatorInstance) ExtraParams(gadgetCtx operators.GadgetContext) api.Params {
	return api.Params{
		{
			Key:          ParamAllowHostFS,
			Description:  "allow access to host filesystem",
			DefaultValue: "false",
			TypeHint:     api.TypeBool,
		},
	}
}

func (i *wasmOperatorInstance) addHandle(obj any) uint32 {
	if obj == nil {
		return 0
	}

	i.handleLock.Lock()
	defer i.handleLock.Unlock()
	i.handleCtr++
	if i.handleCtr == 0 { // 0 is reserved
		i.handleCtr++
	}
	xctr := 0
	for {
		if xctr > 1<<32 {
			// exhausted; TODO: report somehow
			return 0
		}
		if _, ok := i.handleMap[i.handleCtr]; !ok {
			// register new entry
			i.handleMap[i.handleCtr] = obj
			return i.handleCtr
		}
		xctr++
	}
}

func (i *wasmOperatorInstance) getHandle(handleID uint32) any {
	i.handleLock.RLock()
	defer i.handleLock.RUnlock()
	return i.handleMap[handleID]
}

func (i *wasmOperatorInstance) delHandle(handleID uint32) {
	i.handleLock.Lock()
	defer i.handleLock.Unlock()
	delete(i.handleMap, handleID)
}

func (i *wasmOperatorInstance) init(gadgetCtx operators.GadgetContext) error {
	ctx := gadgetCtx.Context()
	i.rt = wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCloseOnContextDone(true))
	// TODO: add mem limits etc

	config := wazero.NewModuleConfig().
		WithStdout(os.Stdout).WithStderr(os.Stderr).WithSysWalltime()

	if i.allowHostFS {
		config = config.WithFS(os.DirFS("/"))
	}

	env := i.rt.NewHostModuleBuilder("env")
	env.NewFunctionBuilder().
		WithGoModuleFunction(wapi.GoModuleFunc(func(ctx context.Context, m wapi.Module, stack []uint64) {
			// TODO: implement log level
			buf, err := stringFromStack(m, stack[0])
			if err != nil {
				gadgetCtx.Logger().Warnf("reading string from stack: %v", err)
				return
			}
			gadgetCtx.Logger().Info(buf)
		}), []wapi.ValueType{wapi.ValueTypeI64}, []wapi.ValueType{}).Export("xlog")

	i.addDataSourceFuncs(env)
	i.addFieldFuncs(env)

	// TODO: do we need this one?
	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(func(ctx context.Context, mod wapi.Module, stack []uint64) {
				free := mod.ExportedFunction("free")
				free.Call(ctx, stack[0])
			}),
			[]wapi.ValueType{wapi.ValueTypeI32}, // ptr
			[]wapi.ValueType{},
		).
		Export("mfree")

	// TODO: do we need this one?
	env.NewFunctionBuilder().
		WithGoModuleFunction(
			wapi.GoModuleFunc(i.free),
			[]wapi.ValueType{wapi.ValueTypeI32}, // any map entry
			[]wapi.ValueType{},
		).
		Export("freeHost")

	env.Instantiate(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, i.rt)

	reader, err := oci.GetContentFromDescriptor(gadgetCtx.Context(), i.desc)
	if err != nil {
		return fmt.Errorf("getting wasm program: %w", err)
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	// TODO: don't pass os.Args[1] directly
	mod, err := i.rt.InstantiateWithConfig(ctx, buf.Bytes(), config.WithArgs("wasi", os.Args[1]))
	if err != nil {
		// Note: Most compilers do not exit the module after running "_start",
		// unless there was an error. This allows you to call exported functions.
		if exitErr, ok := err.(*sys.ExitError); ok && exitErr.ExitCode() != 0 {
			fmt.Fprintf(os.Stderr, "exit_code: %d\n", exitErr.ExitCode())
		} else if !ok {
			log.Panicln(err)
		}
	}
	i.mod = mod

	i.guestMalloc = mod.ExportedFunction("malloc")

	return err
}

func (i *wasmOperatorInstance) Init(gadgetCtx operators.GadgetContext) error {
	fn := i.mod.ExportedFunction("init")
	if fn == nil {
		return nil
	}
	ret, err := fn.Call(gadgetCtx.Context())
	if err != nil {
		return err

	}
	if ret[0] != 0 {
		return errors.New("init failed")
	}
	return nil
}

func (i *wasmOperatorInstance) Start(gadgetCtx operators.GadgetContext) error {
	fn := i.mod.ExportedFunction("start")
	if fn == nil {
		return nil
	}
	ret, err := fn.Call(gadgetCtx.Context())
	if err != nil {
		return err

	}
	if ret[0] != 0 {
		return errors.New("start failed")
	}
	return nil
}

func (i *wasmOperatorInstance) Stop(gadgetCtx operators.GadgetContext) error {
	defer func() {
		// cleanup
		i.handleLock.Lock()
		i.handleMap = nil
		i.handleLock.Unlock()
	}()
	fn := i.mod.ExportedFunction("stop")
	if fn == nil {
		return nil
	}

	// We need a new context in here, as gadgetCtx has already been cancelled
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := fn.Call(ctx)
	return err
}

func init() {
	operators.RegisterOperatorForMediaType(wasmObjectMediaType, &wasmOperator{})
}
