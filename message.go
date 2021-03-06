package pbast

type Message struct {
	Name     string
	Comment  Comment
	Fields   []*MessageField
	Enums    []*Enum
	Messages []*Message
	OneOfs   []*OneOf
}

func NewMessage(name string) *Message {
	return &Message{
		Name: name,
	}
}

func (m *Message) AddField(f *MessageField) *Message {
	if f == nil {
		return m
	}
	m.Fields = append(m.Fields, f)
	return m
}

func (m *Message) AddEnum(e *Enum) *Message {
	if e == nil {
		return m
	}
	m.Enums = append(m.Enums, e)
	return m
}

func (m *Message) AddMessage(n *Message) *Message {
	if n == nil {
		return m
	}
	m.Messages = append(m.Messages, n)
	return m
}

func (m *Message) AddOneOf(o *OneOf) *Message {
	if o == nil {
		return m
	}
	m.OneOfs = append(m.OneOfs, o)
	return m
}

func (m *Message) AddType(t Type) {
	if t == nil {
		return
	}

	switch t := t.(type) {
	case *Message:
		m.AddMessage(t)
	case *Enum:
		m.AddEnum(t)
	default:
		// no-op
	}
}

type MessageField struct {
	Repeated bool
	Type     string
	Name     string
	Index    int
	Options  []*FieldOption
	Comment  Comment
}

type FieldOption struct {
	Name  string
	Value string
}

func NewMessageField(t Type, name string, index int) *MessageField {
	// no repeat by default
	return &MessageField{
		Type:  t.TypeName(),
		Name:  name,
		Index: index,
	}
}

func NewRepeatedMessageField(t Type, name string, index int) *MessageField {
	return &MessageField{
		Repeated: true,
		Type:     t.TypeName(),
		Name:     name,
		Index:    index,
	}
}

func NewFieldOption(name, value string) *FieldOption {
	return &FieldOption{
		Name:  name,
		Value: value,
	}
}

func (f *MessageField) AddOption(o *FieldOption) *MessageField {
	if o == nil {
		return f
	}
	f.Options = append(f.Options, o)
	return f
}
