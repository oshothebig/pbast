package pbast

type Message struct {
	Name     string
	Fields   []MessageField
	Enums    []Enum
	Messages []Message
}

func NewMessage(name string) Message {
	return Message{
		Name: name,
	}
}

func (m Message) AddField(f MessageField) Message {
	nm := Message(m)
	nm.Fields = append(nm.Fields, f)
	return nm
}

func (m Message) AddEnum(e Enum) Message {
	nm := Message(m)
	nm.Enums = append(nm.Enums, e)
	return nm
}

func (m Message) AddMessage(n Message) Message {
	nm := Message(m)
	nm.Messages = append(nm.Messages, n)
	return nm
}

type MessageField struct {
	Repeated bool
	Type     string
	Name     string
	Index    int
	Options  []FieldOption
}

type FieldOption struct {
	Name  string
	Value string
}

func NewMessageField(t Type, name string, index int) MessageField {
	// no repeat by default
	return MessageField{
		Type:  t.TypeName(),
		Name:  name,
		Index: index,
	}
}

func NewRepeatedMessageField(t Type, name string, index int) MessageField {
	return MessageField{
		Repeated: true,
		Type:     t.TypeName(),
		Name:     name,
		Index:    index,
	}
}

func NewFieldOption(name, value string) FieldOption {
	return FieldOption{
		Name:  name,
		Value: value,
	}
}

func (f MessageField) AddOption(o FieldOption) MessageField {
	nf := MessageField(f)
	nf.Options = append(nf.Options, o)
	return nf
}
