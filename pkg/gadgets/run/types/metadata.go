// Copyright 2023 The Inspektor Gadget authors
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

package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"

	"github.com/inspektor-gadget/inspektor-gadget/pkg/columns"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/params"
)

// Keep this aligned with include/gadget/macros.h
const (
	// Prefix used to mark tracer maps
	tracerInfoPrefix = "gadget_tracer_"

	// Prefix used to mark topper maps
	topperInfoPrefix = "gadget_topper_"

	// Prefix used to mark eBPF params
	paramPrefix = "gadget_param_"

	// Prefix used to mark snapshotters structs
	snapshottersPrefix = "gadget_snapshotter_"

	// Prefix used to mark tracer map created with GADGET_TRACER_MAP() defined in
	// include/gadget/buffer.h.
	TracerMapPrefix = "gadget_map_tracer_"
)

// Keep this aligned with include/gadget/types.h
const (
	// Name of the type that gadgets should use to store an L3 endpoint.
	L3EndpointTypeName = "gadget_l3endpoint_t"

	// Name of the type that gadgets should use to store an L4 endpoint.
	L4EndpointTypeName = "gadget_l4endpoint_t"

	// Name of the type to store a mount namespace inode id
	MntNsIdTypeName = "gadget_mntns_id"

	// Name of the type to store a timestamp
	TimestampTypeName = "gadget_timestamp"
)

// Keep this aligned with include/gadget/buffer.h
const (
	GadgetHeapMapName = "gadget_heap"
)

const (
	// Name of the parameter that defins the network interface a TC program is attached to.
	IfaceParam = "iface"
)

type EBPFParam struct {
	params.ParamDesc `yaml:",inline"`
}

const (
	DefaultColumnWidth = 16
)

type Alignment string

const (
	AlignmenNone   Alignment = ""
	AlignmentLeft  Alignment = "left"
	AlignmentRight Alignment = "right"
)

type EllipsisType string

const (
	EllipsisNone   EllipsisType = ""
	EllipsisStart  EllipsisType = "start"
	EllipsisMiddle EllipsisType = "middle"
	EllipsisEnd    EllipsisType = "end"
)

// FieldAttributes describes how to format a field. It's almost 1:1 mapping with columns.Attributes,
// however we are keeping this separated because we don't want to create a strong coupling with the
// columns library now. Later on we can consider merging both of them.
type FieldAttributes struct {
	// Width to reserve for this field
	Width uint `yaml:"width,omitempty"`
	// MinWidth is the minimum width for this field
	MinWidth uint `yaml:"minWidth,omitempty"`
	// MaxWidth is the maximum width for this field
	MaxWidth uint `yaml:"maxWidth,omitempty"`
	// Alignment of this column (left or right)
	Alignment Alignment `yaml:"alignment,omitempty"`
	// Hidden defines whether a column is to be hid by default
	Hidden bool `yaml:"hidden,omitempty"`
	// EllipsisType defines how to abbreviate this column if the value needs more space than is
	// available. (start, middle or end)
	Ellipsis EllipsisType `yaml:"ellipsis,omitempty"`
	// Template defines the template that will be used.
	// TODO: add a link to existing templates
	Template string `yaml:"template,omitempty"`
}

type Field struct {
	// Field name
	Name string `yaml:"name"`
	// Field description
	Description string `yaml:"description,omitempty"`
	// Attributes defines how the field should be formatted
	Attributes FieldAttributes `yaml:"attributes"`
	// Annotations represents extra information that is not relevant to Inspektor Gadget, but
	// for other applications, like color font for instance.
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
}

// Struct describes a type generated by the gadget
type Struct struct {
	Fields []Field `yaml:"fields"`
}

// Tracer describes the behavior of a gadget that collects and sends events to user space
// TODO: We need to rename this concept not to collide with the opentelemetry concept
type Tracer struct {
	// Name of the perf event array or ring buffer that the gadget uses to send events
	MapName string `yaml:"mapName"`
	// Name of the structure generated by this tracer
	StructName string `yaml:"structName"`
}

// Topper describes the behavior of a gadget that shows the current activity
// sorted by the highest to the lowest in the resource being observed.
type Topper struct {
	// Name of the hash map that the gadget uses to send statistics
	MapName string `yaml:"mapName"`
	// Name of the structure generated by this topper
	StructName string `yaml:"structName"`
}

// Snapshotter describes the behavior of a gadget that collects the state of a subsystem
type Snapshotter struct {
	StructName string `yaml:"structName"`
}

// Profiler describes the behavior of a gadget profiling a given subsystem.
type Profiler struct {
	// Name of the structure used as key of the below map.
	StructKeyName string `yaml:"structKeyName"`

	// Name of the structure used as value of the below map.
	StructValueName string `yaml:"structValueName"`

	// Name of the hash map the gadget uses to send events.
	MapName string `yaml:"mapName"`
}

type GadgetMetadata struct {
	// Gadget name
	Name string `yaml:"name"`
	// Gadget description
	Description string `yaml:"description,omitempty"`
	// HomepageURL is the URL to the gadget's homepage
	HomepageURL string `yaml:"homepageURL,omitempty"`
	// DocumentationURL is the URL to the gadget's documentation
	DocumentationURL string `yaml:"documentationURL,omitempty"`
	// SourceURL is the URL to the gadget's source code repository
	SourceURL string `yaml:"sourceURL,omitempty"`
	// Annotations is a map of key-value pairs that provide additional information about the gadget
	Annotations map[string]string `yaml:"annotations,omitempty"`

	// Tracers implemented by the gadget
	// TODO: Rename this field to something that doesn't collide with the opentelemetry concept
	Tracers map[string]Tracer `yaml:"tracers,omitempty"`
	// Toppers implemented by the gadget
	Toppers map[string]Topper `yaml:"toppers,omitempty"`
	// Snapshotters implemented by the gadget
	Snapshotters map[string]Snapshotter `yaml:"snapshotters,omitempty"`
	// Profilers implemented by the gadget
	Profilers map[string]Profiler `yaml:"profilers,omitempty"`
	// Types generated by the gadget
	Structs map[string]Struct `yaml:"structs,omitempty"`
	// Params exposed by the gadget through eBPF constants
	EBPFParams map[string]EBPFParam `yaml:"ebpfParams,omitempty"`
	// Other params exposed by the gadget
	GadgetParams map[string]params.ParamDesc `yaml:"gadgetParams,omitempty"`
}

// countDistImp returns the number of distinct implementations of tracers,
// snapshotters and toppers that the gadget has.
func (m *GadgetMetadata) countDistImp() int {
	count := 0
	if len(m.Tracers) > 0 {
		count++
	}
	if len(m.Snapshotters) > 0 {
		count++
	}
	if len(m.Toppers) > 0 {
		count++
	}
	if len(m.Profilers) > 0 {
		count++
	}
	return count
}

func (m *GadgetMetadata) Validate(spec *ebpf.CollectionSpec) error {
	var result error

	if m.Name == "" {
		result = multierror.Append(result, errors.New("gadget name is required"))
	}

	// Temporary limitation
	if count := m.countDistImp(); count > 1 {
		result = multierror.Append(
			result,
			fmt.Errorf("gadget can implement only one tracer or snapshotter or topper or profiler, found %d", count),
		)
	}

	if err := m.validateEbpfParams(spec); err != nil {
		result = multierror.Append(result, err)
	}

	if err := m.validateTracers(spec); err != nil {
		result = multierror.Append(result, err)
	}

	if err := m.validateToppers(spec); err != nil {
		result = multierror.Append(result, err)
	}

	if err := m.validateSnapshotters(spec); err != nil {
		result = multierror.Append(result, err)
	}

	if err := m.validateProfilers(spec); err != nil {
		result = multierror.Append(result, err)
	}

	if err := m.validateStructs(spec); err != nil {
		result = multierror.Append(result, err)
	}

	if err := m.validateGadgetParams(spec); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

func (m *GadgetMetadata) validateTracers(spec *ebpf.CollectionSpec) error {
	var result error

	// Temporary limitation
	if len(m.Tracers) > 1 {
		result = multierror.Append(result, errors.New("only one tracer is allowed"))
	}

	for name, t := range m.Tracers {
		err := validateMapAndStruct(t.MapName, t.StructName, spec, m, validateTracerMap)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("validating tracer %q: %w", name, err))
		}
	}

	return result
}

// validateTracerMap only checks if the map type. It does not check the map
// value name and type because such a information is not available in the map
// definition for perf event arrays and ring buffers.
func validateTracerMap(tracerMap *ebpf.MapSpec, _ string) error {
	if tracerMap.Type != ebpf.RingBuf && tracerMap.Type != ebpf.PerfEventArray {
		return fmt.Errorf("map %q has a wrong type, expected: ringbuf or perf event array, got: %s",
			tracerMap.Name, tracerMap.Type)
	}
	return nil
}

func (m *GadgetMetadata) validateToppers(spec *ebpf.CollectionSpec) error {
	var result error

	// Temporary limitation
	if len(m.Toppers) > 1 {
		result = multierror.Append(result, errors.New("only one topper is allowed"))
	}

	for name, t := range m.Toppers {
		err := validateMapAndStruct(t.MapName, t.StructName, spec, m, validateTopperMap)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("validating topper %q: %w", name, err))
		}
	}

	return result
}

func validateTopperMap(topperMap *ebpf.MapSpec, expectedStructName string) error {
	if topperMap.Type != ebpf.Hash {
		return fmt.Errorf("map %q has a wrong type, expected: hash, got: %s",
			topperMap.Name, topperMap.Type)
	}

	if topperMap.Value == nil {
		return fmt.Errorf("map %q does not have BTF information for its values", topperMap.Name)
	}

	topperMapStruct, ok := topperMap.Value.(*btf.Struct)
	if !ok {
		return fmt.Errorf("map %q value is %q, expected \"struct\"",
			topperMap.Name, topperMap.Value.TypeName())
	}

	if expectedStructName != "" && topperMapStruct.Name != expectedStructName {
		return fmt.Errorf("map %q value name is %q, expected %q",
			topperMap.Name, topperMapStruct.Name, expectedStructName)
	}

	return nil
}

func (m *GadgetMetadata) validateSnapshotters(spec *ebpf.CollectionSpec) error {
	var result error

	// Temporary limitation
	if len(m.Snapshotters) > 1 {
		result = multierror.Append(result, errors.New("only one snapshotter is allowed"))
	}

	for name, snapshotter := range m.Snapshotters {
		if snapshotter.StructName == "" {
			result = multierror.Append(result, fmt.Errorf("snapshotter %q is missing structName", name))
			continue
		}

		if _, ok := m.Structs[snapshotter.StructName]; !ok {
			result = multierror.Append(result, fmt.Errorf("snapshotter %q references unknown struct %q", name, snapshotter.StructName))
		}
	}

	return result
}

// validateMapAndStruct fully validates the map, while the struct is only
// checked for existence in the Structs section of the metadata as it will be
// validated with the rest of the structs.
func validateMapAndStruct(mapName, structName string,
	spec *ebpf.CollectionSpec,
	m *GadgetMetadata,
	validateMap func(*ebpf.MapSpec, string) error,
) (result error) {
	if mapName == "" {
		result = multierror.Append(result, errors.New("missing mapName"))
	} else {
		ebpfMap, ok := spec.Maps[mapName]
		if !ok {
			return fmt.Errorf("map %q not found in eBPF object", mapName)
		}

		if err := validateMap(ebpfMap, structName); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if structName == "" {
		result = multierror.Append(result, errors.New("missing structName"))
	} else if _, ok := m.Structs[structName]; !ok {
		result = multierror.Append(result, fmt.Errorf("referencing unknown struct %q", structName))
	}

	return
}

func (m *GadgetMetadata) validateProfilers(spec *ebpf.CollectionSpec) error {
	var result error

	// Temporary limitation
	if len(m.Profilers) > 1 {
		result = multierror.Append(result, errors.New("only one profiler is allowed"))
	}

	for name, profiler := range m.Profilers {
		if profiler.StructKeyName == "" {
			result = multierror.Append(result, fmt.Errorf("profiler %q is missing structKeyName", name))
		}

		_, ok := m.Structs[profiler.StructKeyName]
		if !ok {
			result = multierror.Append(result, fmt.Errorf("profiler %q references unknown key struct %q", name, profiler.StructKeyName))
		}

		err := validateMapAndStruct(profiler.MapName, profiler.StructValueName, spec, m, validateProfilerMap)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("validating profiler %q: %w", name, err))
		}
	}

	return result
}

func validateProfilerMap(profilerMap *ebpf.MapSpec, _ string) error {
	if profilerMap.Type != ebpf.Hash {
		return fmt.Errorf("map %q has a wrong type, expected: hash, got: %s", profilerMap.Name, profilerMap.Type)
	}

	return nil
}

func (m *GadgetMetadata) validateStructs(spec *ebpf.CollectionSpec) error {
	var result error

	for name, mapStruct := range m.Structs {
		var btfStruct *btf.Struct
		if err := spec.Types.TypeByName(name, &btfStruct); err != nil {
			result = multierror.Append(result, fmt.Errorf("looking for struct %q in eBPF object: %w", name, err))
			continue
		}

		mapStructFields := make(map[string]Field, len(mapStruct.Fields))
		for _, f := range mapStruct.Fields {
			mapStructFields[f.Name] = f
		}

		btfStructFields := make(map[string]btf.Member, len(btfStruct.Members))
		for _, m := range btfStruct.Members {
			btfStructFields[m.Name] = m
		}

		for fieldName := range mapStructFields {
			if _, ok := btfStructFields[fieldName]; !ok {
				result = multierror.Append(result, fmt.Errorf("field %q not found in eBPF struct %q", fieldName, name))
			}
		}
	}

	return result
}

func (m *GadgetMetadata) validateEbpfParams(spec *ebpf.CollectionSpec) error {
	var result error
	for varName := range m.EBPFParams {
		if err := checkParamVar(spec, varName); err != nil {
			result = multierror.Append(result, err)
		}
		if len(m.EBPFParams[varName].Key) == 0 {
			result = multierror.Append(result, fmt.Errorf("param %q has an empty key", varName))
		}
	}
	return result
}

func (m *GadgetMetadata) validateGadgetParams(spec *ebpf.CollectionSpec) error {
	var result error
	for _, p := range spec.Programs {
		switch p.Type {
		// Networking programs provide an interface name to attach to
		case ebpf.SchedCLS:
			if len(m.GadgetParams) == 0 {
				result = multierror.Append(result, fmt.Errorf("there aren't gadget parameters"))
			} else {
				if _, ok := m.GadgetParams[IfaceParam]; !ok {
					result = multierror.Append(result, fmt.Errorf("iface param not found"))
				}
			}
		}
	}
	return result
}

// Populate fills the metadata from its ebpf spec
func (m *GadgetMetadata) Populate(spec *ebpf.CollectionSpec) error {
	if m.Name == "" {
		m.Name = "TODO: Fill the gadget name"
	}

	if m.Description == "" {
		m.Description = "TODO: Fill the gadget description"
	}

	if m.HomepageURL == "" {
		m.HomepageURL = "TODO: Fill the gadget homepage URL"
	}

	if m.DocumentationURL == "" {
		m.DocumentationURL = "TODO: Fill the gadget documentation URL"
	}

	if m.SourceURL == "" {
		m.SourceURL = "TODO: Fill the gadget source code URL"
	}

	if err := m.populateTracers(spec); err != nil {
		return fmt.Errorf("handling trace maps: %w", err)
	}

	if err := m.populateToppers(spec); err != nil {
		return fmt.Errorf("handling toppers: %w", err)
	}

	if err := m.populateSnapshotters(spec); err != nil {
		return fmt.Errorf("handling snapshotters: %w", err)
	}

	if err := m.populateEbpfParams(spec); err != nil {
		return fmt.Errorf("handling params: %w", err)
	}

	if err := m.populateGadgetParams(spec); err != nil {
		return fmt.Errorf("handling gadget params: %w", err)
	}

	return nil
}

func getUnderlyingType(tf *btf.Typedef) (btf.Type, error) {
	switch typedMember := tf.Type.(type) {
	case *btf.Typedef:
		return getUnderlyingType(typedMember)
	default:
		return typedMember, nil
	}
}

func getColumnSize(typ btf.Type) uint {
	switch typedMember := typ.(type) {
	case *btf.Int:
		switch typedMember.Encoding {
		case btf.Signed:
			switch typedMember.Size {
			case 1:
				return columns.MaxCharsInt8
			case 2:
				return columns.MaxCharsInt16
			case 4:
				return columns.MaxCharsInt32
			case 8:
				return columns.MaxCharsInt64

			}
		case btf.Unsigned:
			switch typedMember.Size {
			case 1:
				return columns.MaxCharsUint8
			case 2:
				return columns.MaxCharsUint16
			case 4:
				return columns.MaxCharsUint32
			case 8:
				return columns.MaxCharsUint64
			}
		case btf.Bool:
			return columns.MaxCharsBool
		case btf.Char:
			return columns.MaxCharsChar
		}
	case *btf.Typedef:
		typ, _ := getUnderlyingType(typedMember)
		return getColumnSize(typ)
	}

	return DefaultColumnWidth
}

func (m *GadgetMetadata) populateTracers(spec *ebpf.CollectionSpec) error {
	tracerInfo, err := getTracerInfo(spec)
	if err != nil {
		return err
	}
	if tracerInfo == nil {
		log.Debug("No tracer found in eBPF object")
		return nil
	}

	if m.Tracers == nil {
		m.Tracers = make(map[string]Tracer)
	}

	tracerMap := spec.Maps[tracerInfo.mapName]
	if tracerMap == nil {
		return fmt.Errorf("map %q not found in eBPF object", tracerInfo.mapName)
	}

	if err := validateTracerMap(tracerMap, ""); err != nil {
		return fmt.Errorf("tracer map is invalid: %w", err)
	}

	var tracerMapStruct *btf.Struct
	if err := spec.Types.TypeByName(tracerInfo.eventType, &tracerMapStruct); err != nil {
		return fmt.Errorf("finding struct %q in eBPF object: %w", tracerInfo.eventType, err)
	}

	if _, found := m.Tracers[tracerInfo.name]; !found {
		log.Debugf("Adding tracer %q with map %q and struct %q",
			tracerInfo.name, tracerMap.Name, tracerMapStruct.Name)

		m.Tracers[tracerInfo.name] = Tracer{
			MapName:    tracerMap.Name,
			StructName: tracerMapStruct.Name,
		}
	} else {
		log.Debugf("Tracer %q already defined, skipping", tracerInfo.name)
	}

	if err := m.populateStruct(tracerMapStruct); err != nil {
		return fmt.Errorf("populating struct: %w", err)
	}

	return nil
}

func (m *GadgetMetadata) populateToppers(spec *ebpf.CollectionSpec) error {
	topperInfo, err := getTopperInfo(spec)
	if err != nil {
		return err
	}
	if topperInfo == nil {
		log.Debug("No topper found in eBPF object")
		return nil
	}

	if m.Toppers == nil {
		m.Toppers = make(map[string]Topper)
	}

	topperMap := spec.Maps[topperInfo.mapName]
	if topperMap == nil {
		return fmt.Errorf("map %q not found in eBPF object", topperInfo.mapName)
	}

	t, found := m.Toppers[topperInfo.name]
	if err := validateTopperMap(topperMap, t.StructName); err != nil {
		return err
	}

	var topperMapStruct *btf.Struct
	if err := spec.Types.TypeByName(topperMap.Value.TypeName(), &topperMapStruct); err != nil {
		return fmt.Errorf("finding struct %q in eBPF object: %w", topperMap.Value.TypeName(), err)
	}

	if !found {
		log.Debugf("Adding topper %q with map %q and struct %q",
			topperInfo.name, topperMap.Name, topperMapStruct.Name)

		m.Toppers[topperInfo.name] = Topper{
			MapName:    topperMap.Name,
			StructName: topperMapStruct.Name,
		}
	} else {
		log.Debugf("Topper %q already defined, skipping", topperInfo.name)
	}

	if err := m.populateStruct(topperMapStruct); err != nil {
		return fmt.Errorf("populating struct: %w", err)
	}

	return nil
}

// GetGadgetIdentByPrefix returns the strings generated by GADGET_ macros.
func GetGadgetIdentByPrefix(spec *ebpf.CollectionSpec, prefix string) ([]string, error) {
	var resultNames []string
	var resultError error

	it := spec.Types.Iterate()
	for it.Next() {
		btfVar, ok := it.Type.(*btf.Var)
		if !ok {
			continue
		}
		if !strings.HasPrefix(btfVar.Name, prefix) {
			continue
		}
		if btfVar.Linkage != btf.GlobalVar {
			resultError = multierror.Append(resultError, fmt.Errorf("%q is not a global variable", btfVar.Name))
		}
		btfPtr, ok := btfVar.Type.(*btf.Pointer)
		if !ok {
			resultError = multierror.Append(resultError, fmt.Errorf("%q is not a pointer", btfVar.Name))
			continue
		}
		btfConst, ok := btfPtr.Target.(*btf.Const)
		if !ok {
			resultError = multierror.Append(resultError, fmt.Errorf("%q is not const", btfVar.Name))
			continue
		}
		_, ok = btfConst.Type.(*btf.Void)
		if !ok {
			resultError = multierror.Append(resultError, fmt.Errorf("%q is not a const void pointer", btfVar.Name))
			continue
		}

		resultNames = append(resultNames, strings.TrimPrefix(btfVar.Name, prefix))
	}

	return resultNames, resultError
}

type tracerInfo struct {
	name      string
	mapName   string
	eventType string
}

// getTracerInfo returns the tracer info generated with GADGET_TRACER().
// If there are multiple annotations only the first one is returned.
func getTracerInfo(spec *ebpf.CollectionSpec) (*tracerInfo, error) {
	tracersInfo, err := GetGadgetIdentByPrefix(spec, tracerInfoPrefix)
	if err != nil {
		return nil, err
	}
	if len(tracersInfo) == 0 {
		return nil, nil
	}

	if len(tracersInfo) > 1 {
		log.Warnf("multiple tracers found, using %q", tracersInfo[0])
	}

	parts := strings.Split(tracersInfo[0], "___")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid tracer info: %q", tracersInfo[0])
	}

	return &tracerInfo{
		name:      parts[0],
		mapName:   parts[1],
		eventType: parts[2],
	}, nil
}

type topperInfo struct {
	name    string
	mapName string
}

// getTopperInfo returns the topper info generated with GADGET_TOPPER().
// If there are multiple annotations only the first one is returned.
func getTopperInfo(spec *ebpf.CollectionSpec) (*topperInfo, error) {
	toppersInfo, err := GetGadgetIdentByPrefix(spec, topperInfoPrefix)
	if err != nil {
		return nil, fmt.Errorf("getting topper info: %w", err)
	}
	if len(toppersInfo) == 0 {
		return nil, nil
	}

	if len(toppersInfo) > 1 {
		log.Warnf("multiple toppers found, using %q", toppersInfo[0])
	}

	parts := strings.Split(toppersInfo[0], "___")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid topper info: %q", toppersInfo[0])
	}

	return &topperInfo{
		name:    parts[0],
		mapName: parts[1],
	}, nil
}

func (m *GadgetMetadata) populateStruct(btfStruct *btf.Struct) error {
	if m.Structs == nil {
		m.Structs = make(map[string]Struct)
	}

	gadgetStruct := m.Structs[btfStruct.Name]
	existingFields := make(map[string]struct{})
	for _, field := range gadgetStruct.Fields {
		existingFields[field.Name] = struct{}{}
	}

	for _, member := range btfStruct.Members {
		// check if field already exists
		if _, ok := existingFields[member.Name]; ok {
			log.Debugf("Field %q already exists, skipping", member.Name)
			continue
		}

		log.Debugf("Adding field %q", member.Name)
		field := Field{
			Name:        member.Name,
			Description: "TODO: Fill field description",
			Attributes: FieldAttributes{
				Width:     getColumnSize(member.Type),
				Alignment: AlignmentLeft,
				Ellipsis:  EllipsisEnd,
			},
		}

		gadgetStruct.Fields = append(gadgetStruct.Fields, field)
	}

	m.Structs[btfStruct.Name] = gadgetStruct

	return nil
}

func (m *GadgetMetadata) populateEbpfParams(spec *ebpf.CollectionSpec) error {
	var result error

	paramNames, err := GetGadgetIdentByPrefix(spec, paramPrefix)
	if err != nil {
		result = multierror.Append(result, err)
	}

	for _, name := range paramNames {
		var btfVar *btf.Var
		err := spec.Types.TypeByName(name, &btfVar)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("looking variable %q up: %w", name, err))
			continue
		}

		err = checkParamVar(spec, name)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		if m.EBPFParams == nil {
			m.EBPFParams = make(map[string]EBPFParam)
		}

		if _, found := m.EBPFParams[name]; found {
			log.Debugf("Param %q already defined, skipping", name)
			continue
		}

		log.Debugf("Adding param %q", name)
		m.EBPFParams[name] = EBPFParam{
			ParamDesc: params.ParamDesc{
				Key:         name,
				Description: "TODO: Fill parameter description",
			},
		}
	}

	return result
}

func (m *GadgetMetadata) populateGadgetParams(spec *ebpf.CollectionSpec) error {
	for _, p := range spec.Programs {
		switch p.Type {
		// Networking programs provide an interface name to attach to
		case ebpf.SchedCLS:
			if m.GadgetParams == nil {
				m.GadgetParams = make(map[string]params.ParamDesc)
			}

			m.GadgetParams[IfaceParam] = params.ParamDesc{
				Key:         IfaceParam,
				Description: "Network interface to attach to",
			}
		}
	}

	return nil
}

func checkParamVar(spec *ebpf.CollectionSpec, name string) error {
	var result error

	var btfVar *btf.Var
	err := spec.Types.TypeByName(name, &btfVar)
	if err != nil {
		result = multierror.Append(result, fmt.Errorf("variable %q not found in eBPF object: %w", name, err))
		return result
	}
	if btfVar.Linkage != btf.GlobalVar {
		result = multierror.Append(result, fmt.Errorf("%q is not a global variable", name))
	}
	btfConst, ok := btfVar.Type.(*btf.Const)
	if !ok {
		result = multierror.Append(result, fmt.Errorf("%q is not const", name))
		return result
	}
	_, ok = btfConst.Type.(*btf.Volatile)
	if !ok {
		result = multierror.Append(result, fmt.Errorf("%q is not volatile", name))
		return result
	}

	return result
}

func (m *GadgetMetadata) populateSnapshotters(spec *ebpf.CollectionSpec) error {
	snapshottersNameAndType, _ := GetGadgetIdentByPrefix(spec, snapshottersPrefix)
	if len(snapshottersNameAndType) == 0 {
		log.Debug("No snapshotters found")
		return nil
	}

	if len(snapshottersNameAndType) > 1 {
		log.Warnf("Multiple snapshotters found, using %q", snapshottersNameAndType[0])
	}

	snapshotterNameAndType := snapshottersNameAndType[0]

	if m.Snapshotters == nil {
		m.Snapshotters = make(map[string]Snapshotter)
	}

	parts := strings.Split(snapshotterNameAndType, "___")
	if len(parts) != 2 {
		return fmt.Errorf("invalid snapshotter annotation: %q", snapshotterNameAndType)
	}
	sname := parts[0]
	stype := parts[1]

	var btfStruct *btf.Struct
	spec.Types.TypeByName(stype, &btfStruct)

	if btfStruct == nil {
		return fmt.Errorf("struct %q not found", stype)
	}

	_, ok := m.Snapshotters[sname]
	if !ok {
		log.Debugf("Adding snapshotter %q", sname)
		m.Snapshotters[sname] = Snapshotter{
			StructName: btfStruct.Name,
		}
	} else {
		log.Debugf("Snapshotter %q already defined, skipping", sname)
	}

	if err := m.populateStruct(btfStruct); err != nil {
		return fmt.Errorf("populating struct: %w", err)
	}

	return nil
}
