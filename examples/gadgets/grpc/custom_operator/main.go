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
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/inspektor-gadget/inspektor-gadget/pkg/datasource"
	igjson "github.com/inspektor-gadget/inspektor-gadget/pkg/datasource/formatters/json"
	gadgetcontext "github.com/inspektor-gadget/inspektor-gadget/pkg/gadget-context"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/operators"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/operators/simple"
	grpcruntime "github.com/inspektor-gadget/inspektor-gadget/pkg/runtime/grpc"
)

func do() error {
	// Define an operator that will be executed on the data received from the server.
	// In this example we implement our own operator, but it's also possible to use existing operators
	// like the cli operator.
	const opPriority = 50000
	myOperator := simple.New("myHandler", simple.OnInit(func(gadgetCtx operators.GadgetContext) error {
		for _, d := range gadgetCtx.GetDataSources() {
			jsonEncoder, _ := igjson.New(d)

			d.Subscribe(func(source datasource.DataSource, data datasource.Data) error {
				jsonOutput := jsonEncoder.Marshal(data)
				fmt.Printf("%s\n", jsonOutput)
				return nil
			}, opPriority)
		}
		return nil
	}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	gadgetCtx := gadgetcontext.New(
		ctx,
		"ghcr.io/inspektor-gadget/gadget/trace_open:latest",
		gadgetcontext.WithDataOperators(myOperator),
	)

	runtime := grpcruntime.New()
	runtimeGlobalParams := runtime.GlobalParamDescs().ToParams()
	runtimeGlobalParams.Get(grpcruntime.ParamRemoteAddress).Set("tcp://127.0.0.1:8888")
	//runtimeGlobalParams.Get(grpcruntime.ParamRemoteAddress).Set("unix:///var/run/ig/ig.socket")

	if err := runtime.Init(runtimeGlobalParams); err != nil {
		return fmt.Errorf("runtime init: %w", err)
	}
	defer runtime.Close()

	params := map[string]string{
		// Capture events coming from the host too
		"operator.LocalManager.host": "true",
	}
	if err := runtime.RunGadget(gadgetCtx, nil, params); err != nil {
		return fmt.Errorf("running gadget: %w", err)
	}

	return nil
}

func main() {
	if err := do(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
		os.Exit(1)
	}
}
