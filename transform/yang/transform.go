package yang

import (
	"fmt"
	"os"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/oshothebig/pbast"
)

var builtinMap = map[yang.TypeKind]pbast.Type{
	yang.Yint8:   pbast.Int32,
	yang.Yint16:  pbast.Int32,
	yang.Yint32:  pbast.Int32,
	yang.Yint64:  pbast.Int64,
	yang.Yuint8:  pbast.UInt32,
	yang.Yuint16: pbast.UInt32,
	yang.Yuint32: pbast.UInt32,
	yang.Yuint64: pbast.UInt64,
	yang.Ystring: pbast.String,
	yang.Ybool:   pbast.Bool,
	yang.Ybinary: pbast.Bytes,
}

var builtinTypes = stringSet{
	"int8":                struct{}{},
	"int16":               struct{}{},
	"int32":               struct{}{},
	"int64":               struct{}{},
	"uint8":               struct{}{},
	"uint16":              struct{}{},
	"unit32":              struct{}{},
	"uint64":              struct{}{},
	"string":              struct{}{},
	"boolean":             struct{}{},
	"enumeration":         struct{}{},
	"bits":                struct{}{},
	"binary":              struct{}{},
	"leafref":             struct{}{},
	"identityref":         struct{}{},
	"empty":               struct{}{},
	"union":               struct{}{},
	"instance-identifier": struct{}{},
}

type transformer struct {
	topScope           *scope
	decimal64          *pbast.Message
	leafRef            *pbast.Message
	identityRef        *pbast.Message
	instanceIdentifier *pbast.Message
	emptyNeeded        bool
}

// e must be YANG module
func Transform(e *yang.Entry) *pbast.File {
	if _, ok := e.Node.(*yang.Module); !ok {
		return nil
	}

	t := &transformer{
		topScope: newScope(),
	}

	return t.module(entry{e})
}

func (t *transformer) declare(m pbast.Type) {
	if err := t.topScope.addType(m); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (t *transformer) reflectTo(f *pbast.File) {
	if typ := t.topScope.types[pbast.Empty.TypeName()]; typ != nil {
		f.AddImport(pbast.NewImport("google/protobuf/empty.proto"))
	}
	t.topScope.reflectTo(f)
}

func (t *transformer) module(e entry) *pbast.File {
	namespace := e.Namespace().Name
	f := pbast.NewFile(pbast.NewPackageWithElements(guessElements(namespace)))

	f.Comment = e.moduleComment()

	root := t.buildMessage("Root", e)
	// Clear Root messgage comment because it overlaps with
	// the file level comment being generated from module description too
	root.Comment = nil
	// Child nodes are enclosed with Root message
	f.AddMessage(root)

	// RPCs
	s := t.rpcs(e)
	f.AddService(s)

	// Notifications
	n := t.notifications(e)
	f.AddService(n)

	t.reflectTo(f)

	return f
}

func (t *transformer) rpcs(e entry) *pbast.Service {
	rpcs := e.rpcs()
	if len(rpcs) == 0 {
		return nil
	}

	s := pbast.NewService(CamelCase(e.Name) + "Service")
	for _, rpc := range rpcs {
		r := t.rpc(rpc)
		s.AddRPC(r)
	}

	return s
}

func (t *transformer) rpc(e entry) *pbast.RPC {
	method := CamelCase(e.Name)
	var in pbast.Type
	if e.RPC.Input == nil {
		in = pbast.Empty
	} else {
		in = t.buildMessage(method+"Request", entry{e.RPC.Input})
	}
	t.declare(in)

	var out pbast.Type
	if e.RPC.Output == nil {
		out = pbast.Empty
	} else {
		out = t.buildMessage(method+"Response", entry{e.RPC.Output})
	}
	t.declare(out)

	rpc := pbast.NewRPC(
		method,
		pbast.NewReturnType(in.TypeName()),
		pbast.NewReturnType(out.TypeName()),
	)
	rpc.Comment = e.genericComments()

	return rpc
}

func (t *transformer) notifications(e entry) *pbast.Service {
	notifications := e.notifications()
	if len(notifications) == 0 {
		return nil
	}

	s := pbast.NewService(CamelCase(e.Name + "Notification"))
	for _, notification := range notifications {
		n := t.notification(notification)
		n.Comment = notification.genericComments()
		s.AddRPC(n)
	}

	return s
}

func (t *transformer) notification(e entry) *pbast.RPC {
	const common = "Notification"
	method := CamelCase(e.Name)
	in := pbast.Empty
	t.declare(in)
	out := t.buildMessage(buildName(method, common, "Response"), e)
	t.declare(out)

	returnType := pbast.NewReturnType(out.TypeName())
	returnType.Streamable = true
	return pbast.NewRPC(method, pbast.NewReturnType(in.TypeName()), returnType)
}

func (t *transformer) buildMessage(name string, e entry) *pbast.Message {
	if e.Entry == nil {
		return nil
	}

	msg := &pbast.Message{
		Name:    name,
		Comment: e.genericComments(),
	}
	scope := newScope()
	for index, child := range e.children() {
		fieldNum := index + 1
		var field *pbast.MessageField
		switch {
		case child.IsLeaf(), child.IsLeafList():
			typ := t.leaf(child.Type, messageName(child))

			if needsInTopScope(typ) {
				t.declare(typ)
			} else {
				if err := scope.addType(typ); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}

			field = &pbast.MessageField{
				Repeated: child.IsLeafList(),
				Type:     typ.TypeName(),
				Name:     underscoreCase(child.Name),
				Index:    fieldNum,
				Comment:  child.genericComments(),
			}
		case child.IsContainer(), child.IsList(), child.IsChoice(), child.IsCase():
			field = t.directory(scope, child, fieldNum, child.IsList())
		}
		msg.AddField(field)
	}

	scope.reflectTo(msg)

	return msg
}

func (t *transformer) leaf(typ *yang.YangType, name string) pbast.Type {
	if needsCustomType(typ.Kind) {
		return t.customType(typ, CamelCase(name))
	}

	if isBuiltinType(typ.Name) {
		return t.mapBuiltinType(typ.Kind)
	}

	inner := t.mapBuiltinType(typ.Kind)
	msg := pbast.NewMessage(CamelCase(name)).
		AddField(pbast.NewMessageField(inner, "value", 1))
	if needsInTopScope(inner) {
		t.declare(inner)
	} else {
		msg.AddType(inner)
	}
	return msg
}

func messageName(e entry) string {
	if isBuiltinType(e.Type.Name) {
		return e.Name
	}

	return e.Type.Name
}

func needsInTopScope(t pbast.Type) bool {
	switch t {
	case decimal64, leafRef, identityRef, instanceIdentifier, pbast.Empty:
		return true
	default:
		return false
	}
}

func (t *transformer) customType(typ *yang.YangType, name string) pbast.Type {
	switch typ.Kind {
	case yang.Yenum:
		return t.customEnum(name, typ.Enum)
	case yang.Ybits:
		return t.customBits(name, typ.Bit)
	case yang.Yunion:
		return t.customUnion(name, typ.Type)
	default:
		return nil
	}
}

func isBuiltinType(name string) bool {
	if _, ok := yang.TypeKindFromName[name]; !ok {
		return false
	}
	return true
}

func needsCustomType(kind yang.TypeKind) bool {
	switch kind {
	case yang.Yenum, yang.Ybits, yang.Yunion:
		return true
	default:
		return false
	}
}

func (t *transformer) mapBuiltinType(kind yang.TypeKind) pbast.Type {
	switch kind {
	case yang.Yint8, yang.Yint16, yang.Yint32:
		return pbast.Int32
	case yang.Yint64:
		return pbast.Int64
	case yang.Yuint8, yang.Yuint16, yang.Yuint32:
		return pbast.UInt32
	case yang.Yuint64:
		return pbast.UInt64
	case yang.Ystring:
		return pbast.String
	case yang.Ybool:
		return pbast.Bool
	case yang.Ybinary:
		return pbast.Bytes
	case yang.Yleafref:
		return leafRef
	case yang.Yidentityref:
		return identityRef
	case yang.Yempty:
		return pbast.Empty
	case yang.YinstanceIdentifier:
		return instanceIdentifier
	case yang.Ydecimal64:
		return decimal64
	default:
		return nil
	}
}

func (t *transformer) customBits(name string, bits *yang.EnumType) *pbast.Message {
	msg := pbast.NewMessage(name)
	for i, n := range bits.Names() {
		v := 1 << uint(bits.Values()[i])
		msg.AddField(pbast.NewMessageField(pbast.Bool, n, v))
	}

	return msg
}

func (t *transformer) customEnum(name string, e *yang.EnumType) *pbast.Message {
	enum := pbast.NewEnum("Value")
	for _, v := range e.Values() {
		n := e.ValueMap()[v]
		enum.AddField(pbast.NewEnumField(constantName(n), int(v)))
	}
	msg := pbast.NewMessage(name).
		AddEnum(enum).
		AddField(pbast.NewMessageField(enum, "value", 1))

	return msg
}

func (t *transformer) customUnion(name string, types []*yang.YangType) *pbast.Message {
	scope := newScope()
	pbTypes := t.unionFields(types, nil, scope)

	oneof := pbast.NewOneOf("value")
	for i, typ := range pbTypes {
		oneof.AddField(pbast.NewOneOfField(typ, underscoreCase(typ.TypeName()), i+1))
	}

	msg := pbast.NewMessage(name).AddOneOf(oneof)
	scope.reflectTo(msg)

	return msg
}

func (t *transformer) unionFields(types []*yang.YangType, pbTypes []pbast.Type, scope *scope) []pbast.Type {
	for _, typ := range types {
		if typ.Kind == yang.Yunion {
			pbTypes = t.unionFields(typ.Type, pbTypes, scope)
			continue
		}

		pbtype := t.leaf(typ, typ.Name)
		if pbtype == nil {
			continue
		}

		pbTypes = append(pbTypes, pbtype)
		scope.addType(pbtype)
	}

	return pbTypes
}

func (t *transformer) directory(scope *scope, e entry, index int, repeated bool) *pbast.MessageField {
	fieldName := underscoreCase(e.Name)
	typeName := CamelCase(e.Name)

	inner := t.buildMessage(typeName, e)
	field := &pbast.MessageField{
		Repeated: repeated,
		Type:     inner.TypeName(),
		Name:     fieldName,
		Index:    index,
	}

	if err := scope.addType(inner); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return field
}
