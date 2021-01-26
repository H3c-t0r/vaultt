// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flatbuf

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type FixedSizeBinary struct {
	_tab flatbuffers.Table
}

func GetRootAsFixedSizeBinary(buf []byte, offset flatbuffers.UOffsetT) *FixedSizeBinary {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &FixedSizeBinary{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *FixedSizeBinary) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *FixedSizeBinary) Table() flatbuffers.Table {
	return rcv._tab
}

/// Number of bytes per value
func (rcv *FixedSizeBinary) ByteWidth() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

/// Number of bytes per value
func (rcv *FixedSizeBinary) MutateByteWidth(n int32) bool {
	return rcv._tab.MutateInt32Slot(4, n)
}

func FixedSizeBinaryStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func FixedSizeBinaryAddByteWidth(builder *flatbuffers.Builder, byteWidth int32) {
	builder.PrependInt32Slot(0, byteWidth, 0)
}
func FixedSizeBinaryEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
