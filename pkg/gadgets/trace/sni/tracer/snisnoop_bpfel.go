// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || amd64p32 || arm || arm64 || mips64le || mips64p32le || mipsle || ppc64le || riscv64
// +build 386 amd64 amd64p32 arm arm64 mips64le mips64p32le mipsle ppc64le riscv64

package tracer

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type snisnoopEndpointsKey struct {
	Key struct {
		Prefixlen uint32
		Data      [0]uint8
	}
	Ip uint32
}

type snisnoopEventT struct{ Name [128]uint8 }

// loadSnisnoop returns the embedded CollectionSpec for snisnoop.
func loadSnisnoop() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_SnisnoopBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load snisnoop: %w", err)
	}

	return spec, err
}

// loadSnisnoopObjects loads snisnoop and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*snisnoopObjects
//	*snisnoopPrograms
//	*snisnoopMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadSnisnoopObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadSnisnoop()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// snisnoopSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type snisnoopSpecs struct {
	snisnoopProgramSpecs
	snisnoopMapSpecs
}

// snisnoopSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type snisnoopProgramSpecs struct {
	IgTraceSni *ebpf.ProgramSpec `ebpf:"ig_trace_sni"`
}

// snisnoopMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type snisnoopMapSpecs struct {
	Endpoints *ebpf.MapSpec `ebpf:"endpoints"`
	Events    *ebpf.MapSpec `ebpf:"events"`
}

// snisnoopObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadSnisnoopObjects or ebpf.CollectionSpec.LoadAndAssign.
type snisnoopObjects struct {
	snisnoopPrograms
	snisnoopMaps
}

func (o *snisnoopObjects) Close() error {
	return _SnisnoopClose(
		&o.snisnoopPrograms,
		&o.snisnoopMaps,
	)
}

// snisnoopMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadSnisnoopObjects or ebpf.CollectionSpec.LoadAndAssign.
type snisnoopMaps struct {
	Endpoints *ebpf.Map `ebpf:"endpoints"`
	Events    *ebpf.Map `ebpf:"events"`
}

func (m *snisnoopMaps) Close() error {
	return _SnisnoopClose(
		m.Endpoints,
		m.Events,
	)
}

// snisnoopPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadSnisnoopObjects or ebpf.CollectionSpec.LoadAndAssign.
type snisnoopPrograms struct {
	IgTraceSni *ebpf.Program `ebpf:"ig_trace_sni"`
}

func (p *snisnoopPrograms) Close() error {
	return _SnisnoopClose(
		p.IgTraceSni,
	)
}

func _SnisnoopClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed snisnoop_bpfel.o
var _SnisnoopBytes []byte
