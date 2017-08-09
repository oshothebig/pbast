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

type WellKnownType string

func (t WellKnownType) TypeName() string {
	return string(t)
}

const (
	Any             WellKnownType = "google.protobuf.Any"
	Api             WellKnownType = "google.protobuf.Api"
	BoolValue       WellKnownType = "google.protobuf.BoolValue"
	BytesValue      WellKnownType = "google.protobuf.BytesValue"
	DoubleValue     WellKnownType = "google.protobuf.DoubleValue"
	Duration        WellKnownType = "google.protobuf.Duration"
	Empty           WellKnownType = "google.protobuf.Empty"
	WellKnownEnum   WellKnownType = "google.protobuf.Enum"
	EnumValue       WellKnownType = "google.protobuf.EnumValue"
	WellKnownField  WellKnownType = "google.protobuf.Field"
	Cardinality     WellKnownType = "google.protobuf.Cardinality"
	Kind            WellKnownType = "google.protobuf.Kind"
	FieldMask       WellKnownType = "google.protobuf.FieldMask"
	FloatValue      WellKnownType = "google.protobuf.FloatValue"
	Int32Value      WellKnownType = "google.protobuf.Int32Value"
	Int64Value      WellKnownType = "google.protobuf.Int64Value"
	ListValue       WellKnownType = "google.protobuf.ListValue"
	Method          WellKnownType = "google.protobuf.Method"
	Mixin           WellKnownType = "google.protobuf.Mixin"
	NullValue       WellKnownType = "google.protobuf.NullValue"
	WellKnownOption WellKnownType = "google.protobuf.Option"
	SourceContext   WellKnownType = "google.protobuf.SourceContext"
	StringValue     WellKnownType = "google.protobuf.StringValue"
	WellKnownStruct WellKnownType = "google.protobuf.Struct"
	WellKnownSyntax WellKnownType = "google.protobuf.Syntax"
	Timestamp       WellKnownType = "google.protobuf.Timestamp"
	UInt32Value     WellKnownType = "google.protobuf.UInt32Value"
	UInt64Value     WellKnownType = "google.protobuf.UInt64Value"
	Value           WellKnownType = "google.protobuf.Value"

// Define well know type for "Type" when coming up with how to resolve naming conflict
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

type TypeAdder interface {
	AddType(Type)
}
