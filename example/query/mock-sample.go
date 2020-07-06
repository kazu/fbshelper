package vfs_schema

import (
	flatbuffers "github.com/google/flatbuffers/go"
	base "github.com/kazu/fbshelper/query/base"
	//"reflect"
)

// .Fields.Type

type Field_0_Type int32

const Field_0_i = 0

func (node FbsRoot) SetVersion(v Field_0_Type) error {

	buf := node.ValueNormal(Field_0_i)
	if len(buf) < base.SizeOfint32 {
		return base.ERR_MORE_BUFFER
	}

	flatbuffers.WriteInt32(buf, int32(v)) // camel .Field.[i].Type

	return nil
}
