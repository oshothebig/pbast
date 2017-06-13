package pbast

import (
	"reflect"
)

type Type interface {
	TypeName() string
}

type BuiltinType string

func (t BuiltinType) TypeName() string {
	return string(t)
}

const (
	Double   BuiltinType = "double"
	Float    BuiltinType = "float"
	Int32    BuiltinType = "int32"
	Int64    BuiltinType = "int64"
	UInt32   BuiltinType = "uint32"
	UInt64   BuiltinType = "uint64"
	SInt32   BuiltinType = "sint32"
	SInt64   BuiltinType = "sint64"
	Fixed32  BuiltinType = "fixed32"
	Fixed64  BuiltinType = "fixed64"
	SFixed32 BuiltinType = "sfixed32"
	SFixed64 BuiltinType = "sfixed64"
	Bool     BuiltinType = "bool"
	String   BuiltinType = "string"
	Bytes    BuiltinType = "bytes"
)

func (e Enum) TypeName() string {
	return e.Name
}

func (e Message) TypeName() string {
	return e.Name
}

func IsSameType(t1, t2 Type) bool {
	if t1.TypeName() != t2.TypeName() {
		return false
	}
	return reflect.DeepEqual(t1, t2)
}
