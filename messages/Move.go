// automatically generated by the FlatBuffers compiler, do not modify

package messages

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Move struct {
	_tab flatbuffers.Struct
}

func (rcv *Move) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Move) R() int8 {
	return rcv._tab.GetInt8(rcv._tab.Pos + flatbuffers.UOffsetT(0))
}
func (rcv *Move) MutateR(n int8) bool {
	return rcv._tab.MutateInt8(rcv._tab.Pos+flatbuffers.UOffsetT(0), n)
}

func (rcv *Move) C() int8 {
	return rcv._tab.GetInt8(rcv._tab.Pos + flatbuffers.UOffsetT(1))
}
func (rcv *Move) MutateC(n int8) bool {
	return rcv._tab.MutateInt8(rcv._tab.Pos+flatbuffers.UOffsetT(1), n)
}

func CreateMove(builder *flatbuffers.Builder, R int8, C int8) flatbuffers.UOffsetT {
	builder.Prep(1, 2)
	builder.PrependInt8(C)
	builder.PrependInt8(R)
	return builder.Offset()
}