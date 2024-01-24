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

//go:build 386 || amd64 || amd64p32 || arm || arm64 || loong64 || mipsle || mips64le || mips64p32le || ppc64le || riscv64

package datasource

import "encoding/binary"

// NativeEndian is set to either binary.BigEndian or binary.LittleEndian,
// depending on the host's endianness. Borrowed from cilium/ebpf.
var NativeEndian = binary.LittleEndian
