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

package execsnoop

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/kinvolk/inspektor-gadget/pkg/gadgets"
	"github.com/kinvolk/inspektor-gadget/pkg/gadgets/execsnoop/tracer"

	coretracer "github.com/kinvolk/inspektor-gadget/pkg/gadgets/execsnoop/tracer/core"
	standardtracer "github.com/kinvolk/inspektor-gadget/pkg/gadgets/execsnoop/tracer/standard"
	"github.com/kinvolk/inspektor-gadget/pkg/gadgets/execsnoop/types"

	gadgetv1alpha1 "github.com/kinvolk/inspektor-gadget/pkg/apis/gadget/v1alpha1"
)

type Trace struct {
	resolver  gadgets.Resolver
	publisher gadgets.FuncPublisher

	started bool
	tracer  tracer.Tracer
}

type TraceFactory struct {
	gadgets.BaseFactory
}

func NewFactory() gadgets.TraceFactory {
	return &TraceFactory{
		BaseFactory: gadgets.BaseFactory{DeleteTrace: deleteTrace},
	}
}

func (f *TraceFactory) Description() string {
	return `execsnoop shows new created processes, with container details.`
}

func (f *TraceFactory) OutputModesSupported() map[string]struct{} {
	return map[string]struct{}{
		"Stream": {},
	}
}

func deleteTrace(name string, t interface{}) {
	trace := t.(*Trace)
	if trace.tracer != nil {
		trace.tracer.Stop()
	}
}

func (f *TraceFactory) Operations() map[string]gadgets.TraceOperation {
	n := func() interface{} {
		return &Trace{
			resolver:  f.Resolver,
			publisher: f.Publisher,
		}
	}

	return map[string]gadgets.TraceOperation{
		"start": {
			Doc: "Start execsnoop gadget",
			Operation: func(name string, trace *gadgetv1alpha1.Trace) {
				f.LookupOrCreate(name, n).(*Trace).Start(trace)
			},
		},
		"stop": {
			Doc: "Stop execsnoop gadget",
			Operation: func(name string, trace *gadgetv1alpha1.Trace) {
				f.LookupOrCreate(name, n).(*Trace).Stop(trace)
			},
		},
	}
}

func (t *Trace) Start(trace *gadgetv1alpha1.Trace) {
	if t.started {
		gadgets.CleanupTraceStatus(trace)
		trace.Status.State = "Started"
		return
	}

	traceName := gadgets.TraceName(trace.ObjectMeta.Namespace, trace.ObjectMeta.Name)

	eventCallback := func(event types.Event) {
		r, err := json.Marshal(event)
		if err != nil {
			log.Warnf("Gadget %s: error marshalling event: %s", trace.Spec.Gadget, err)
			return
		}
		t.resolver.PublishEvent(traceName, string(r))
		t.publisher(trace, event)
	}

	var err error

	config := &tracer.Config{
		MountnsMap: gadgets.TracePinPath(trace.ObjectMeta.Namespace, trace.ObjectMeta.Name),
	}
	t.tracer, err = coretracer.NewTracer(config, t.resolver, eventCallback, trace.Spec.Node)
	if err != nil {
		// TODO: The following line causes the client to fail.
		// trace.Status.OperationWarning = fmt.Sprint("failed to create core tracer. Falling back to standard one")

		// fallback to standard tracer
		log.Infof("Gadget %s: falling back to standard tracer. CO-RE tracer failed: %s",
			trace.Spec.Gadget, err)

		t.tracer, err = standardtracer.NewTracer(config, t.resolver, eventCallback, trace.Spec.Node)
		if err != nil {
			trace.Status.OperationError = fmt.Sprintf("failed to create tracer: %s", err)
			return
		}
	}

	t.started = true

	gadgets.CleanupTraceStatus(trace)
	trace.Status.State = "Started"
}

func (t *Trace) Stop(trace *gadgetv1alpha1.Trace) {
	if !t.started {
		trace.Status.OperationError = "Not started"
		return
	}

	t.tracer.Stop()
	t.tracer = nil
	t.started = false

	trace.Status.OperationError = ""
	trace.Status.State = "Stopped"
}
