package yang

import (
	"fmt"
	"os"
	"strings"

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

type transformer struct {
	topScope  *scope
	decimal64 *pbast.Message
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

func (t *transformer) declare(m *pbast.Message) {
	if m == nil {
		return
	}
	if err := t.topScope.addType(m); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (t *transformer) module(e entry) *pbast.File {
	namespace := e.Namespace().Name
	f := pbast.NewFile(pbast.NewPackageWithElements(guessElements(namespace)))

	f.Comment = t.moduleComment(e)

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

	for _, m := range t.topScope.allTypes() {
		f.AddType(m)
	}
	f.AddMessage(t.decimal64)

	return f
}

func (t *transformer) moduleComment(e entry) pbast.Comment {
	description := t.description(e)
	namespace := t.namespace(e)
	revisions := t.revisions(e)
	reference := t.reference(e)

	var comment []string
	comment = append(comment, description...)
	comment = append(comment, namespace...)
	comment = append(comment, revisions...)
	comment = append(comment, reference...)

	return comment
}

func (t *transformer) genericComments(e entry) pbast.Comment {
	description := t.description(e)
	reference := t.reference(e)

	comments := append(description, reference...)
	return comments
}

func (t *transformer) description(e entry) pbast.Comment {
	description := e.Description
	if e.Description == "" {
		return nil
	}

	lines := strings.Split(strings.TrimRight(description, "\n "), "\n")

	ret := make([]string, 0, len(lines)+1)
	ret = append(ret, "Description:")
	ret = append(ret, lines...)
	return ret
}

func (t *transformer) revisions(e entry) pbast.Comment {
	var lines []string
	if v := e.Extra["revision"]; len(v) > 0 {
		for _, rev := range v[0].([]*yang.Revision) {
			lines = append(lines, "Revision: "+rev.Name)
		}
	}

	return lines
}

func (t *transformer) namespace(e entry) pbast.Comment {
	namespace := e.Namespace().Name
	if namespace == "" {
		return nil
	}

	return []string{"Namespace: " + namespace}
}

func (t *transformer) reference(e entry) pbast.Comment {
	v := e.Extra["reference"]
	if len(v) == 0 {
		return nil
	}

	ref := v[0].(*yang.Value)
	if ref == nil {
		return nil
	}
	if ref.Name == "" {
		return nil
	}

	lines := strings.Split(strings.TrimRight(ref.Name, "\n "), "\n")

	ret := make([]string, 0, len(lines)+1)
	ret = append(ret, "Reference:")
	ret = append(ret, lines...)
	return ret
}

func (t *transformer) rpcs(e entry) *pbast.Service {
	rpcs := e.rpcs()
	if len(rpcs) == 0 {
		return nil
	}

	s := pbast.NewService(CamelCase(e.Name))
	for _, rpc := range rpcs {
		r := t.rpc(rpc)
		s.AddRPC(r)
	}

	return s
}

func (t *transformer) rpc(e entry) *pbast.RPC {
	method := CamelCase(e.Name)
	in := method + "Request"
	out := method + "Response"

	rpc := pbast.NewRPC(
		method,
		pbast.NewReturnType(in),
		pbast.NewReturnType(out),
	)
	rpc.Comment = t.genericComments(e)

	t.declare(t.buildMessage(in, entry{e.RPC.Input}))
	t.declare(t.buildMessage(out, entry{e.RPC.Output}))

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
		n.Comment = t.genericComments(notification)
		s.AddRPC(n)
	}

	return s
}

func (t *transformer) notification(e entry) *pbast.RPC {
	method := CamelCase(e.Name)
	in := method + "NotificationRequest"
	out := method + "NotificationResponse"

	rpc := pbast.NewRPC(method, pbast.NewReturnType(in), pbast.NewReturnType(out))

	// notification statement doesn't have an input parameter equivalent,
	// then empty message is used for input as RPC
	t.declare(pbast.NewMessage(in))
	t.declare(t.buildMessage(out, e))

	return rpc
}

func (t *transformer) buildMessage(name string, e entry) *pbast.Message {
	if e.Entry == nil {
		return nil
	}

	msg := pbast.NewMessage(name)
	msg.Comment = t.genericComments(e)
	for index, child := range e.children() {
		fieldNum := index + 1
		var field *pbast.MessageField
		var inner pbast.Type
		switch {
		// leaf-list case
		case child.Type != nil && child.ListAttr != nil:
			field, inner = t.leaf(child, fieldNum, true)
		// leaf case
		case child.Type != nil:
			field, inner = t.leaf(child, fieldNum, false)
		// list case
		case child.ListAttr != nil:
			inner, field = t.directory(child, fieldNum, true)
		// others might be container case
		default:
			inner, field = t.directory(child, fieldNum, false)
		}
		msg.AddType(inner).AddField(field)
	}

	return msg
}

func (t *transformer) leaf(e entry, index int, repeated bool) (field *pbast.MessageField, nested pbast.Type) {
	typ := builtinMap[e.Type.Kind]
	// no direct builtin type mapping
	// custom message is built
	if typ == nil {
		name := CamelCase(e.Name)
		switch e.Type.Kind {
		// define at the top level
		case yang.Ydecimal64:
			t.decimal64 = decimal64Message
			typ = decimal64Message
		// define as a nested type
		case yang.Ybits:
			typ = t.customBits(name, e.Type.Bit)
		// define as a nested type
		case yang.Yenum:
			typ = t.customEnum(name, e.Type.Enum)
		// not implemented
		case yang.Yunion, yang.Yempty, yang.Yleafref,
			yang.Yidentityref, yang.YinstanceIdentifier:
			return nil, nil
		}
	}

	name := underscoreCase(e.Name)
	field = &pbast.MessageField{
		Repeated: repeated,
		Type:     typ.TypeName(),
		Name:     name,
		Index:    index,
		Comment:  t.genericComments(e),
	}

	if e.Type.Kind == yang.Ydecimal64 {
		return field, nil
	}

	return field, typ
}

func (t *transformer) customBits(name string, bits *yang.EnumType) *pbast.Message {
	msg := pbast.NewMessage(name)
	for i, n := range bits.Names() {
		v := 1 << uint(bits.Values()[i])
		msg.AddField(pbast.NewMessageField(pbast.Bool, n, v))
	}

	return msg
}

func (t *transformer) customEnum(name string, e *yang.EnumType) *pbast.Enum {
	enum := pbast.NewEnum(name)
	for i, n := range e.Names() {
		v := int(e.Values()[i])
		enum.AddField(pbast.NewEnumField(constantName(n), v))
	}

	return enum
}

func (t *transformer) directory(e entry, index int, repeated bool) (*pbast.Message, *pbast.MessageField) {
	fieldName := underscoreCase(e.Name)
	typeName := CamelCase(e.Name)

	inner := t.buildMessage(typeName, e)
	field := &pbast.MessageField{
		Repeated: repeated,
		Type:     inner.TypeName(),
		Name:     fieldName,
		Index:    index,
		Comment:  t.genericComments(e),
	}

	return inner, field
}
