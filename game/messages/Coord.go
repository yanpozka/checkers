// automatically generated by the FlatBuffers compiler, do not modify

package messages

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Coord struct {
	_tab flatbuffers.Struct
}

func (rcv *Coord) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Coord) R() int8 {
	return rcv._tab.GetInt8(rcv._tab.Pos + flatbuffers.UOffsetT(0))
}
func (rcv *Coord) MutateR(n int8) bool {
	return rcv._tab.MutateInt8(rcv._tab.Pos+flatbuffers.UOffsetT(0), n)
}

func (rcv *Coord) C() int8 {
	return rcv._tab.GetInt8(rcv._tab.Pos + flatbuffers.UOffsetT(1))
}
func (rcv *Coord) MutateC(n int8) bool {
	return rcv._tab.MutateInt8(rcv._tab.Pos+flatbuffers.UOffsetT(1), n)
}

func CreateCoord(builder *flatbuffers.Builder, R int8, C int8) flatbuffers.UOffsetT {
	builder.Prep(1, 2)
	builder.PrependInt8(C)
	builder.PrependInt8(R)
	return builder.Offset()
}