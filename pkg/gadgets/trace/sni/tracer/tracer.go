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

package tracer

import (
	"fmt"

	"github.com/inspektor-gadget/inspektor-gadget/pkg/gadgets"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/gadgets/internal/networktracer"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/gadgets/trace/sni/types"
	eventtypes "github.com/inspektor-gadget/inspektor-gadget/pkg/types"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/endpoint-collection"
)

//go:generate bash -c "source ./clangosflags.sh; go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel -cc clang -type event_t snisnoop ./bpf/snisnoop.c -- $CLANG_OS_FLAGS -I./bpf/ -I../../../../endpoint-collection/bpf"

const (
	BPFProgName         = "ig_trace_sni"
	BPFPerfMapName      = "events"
	BPFSocketAttach     = 50
	TLSMaxServerNameLen = len(snisnoopEventT{}.Name)
)

type Config struct {
	EndpointCollection *endpointcollection.EndpointCollection
}

type Tracer struct {
	config *Config

	*networktracer.Tracer[types.Event]
}

func NewTracer(config *Config) (*Tracer, error) {
	spec, err := loadSnisnoop()
	if err != nil {
		return nil, fmt.Errorf("failed to load asset: %w", err)
	}

	return &Tracer{
		config: config,
		Tracer: networktracer.NewTracer(
			spec,
			config.EndpointCollection,
			BPFProgName,
			BPFPerfMapName,
			BPFSocketAttach,
			types.Base,
			parseSNIEvent,
		),
	}, nil
}

func parseSNIEvent(sample []byte) (*types.Event, error) {
	if len(sample) > TLSMaxServerNameLen {
		sample = sample[:TLSMaxServerNameLen]
	}

	name := gadgets.FromCString(sample)
	if len(name) == 0 {
		return nil, nil
	}

	event := types.Event{
		Event: eventtypes.Event{
			Type: eventtypes.NORMAL,
		},
		Name: name,
	}

	return &event, nil
}
