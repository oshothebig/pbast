package pbast

// Node represents a node in an abstract syntax tree
type Node interface {
	name() string
}

func (f *File) name() string {
	return "file"
}

func (s Syntax) name() string {
	return "syntax"
}

func (i *Import) name() string {
	return "import"
}

func (p Package) name() string {
	return "package"
}

func (p *Option) name() string {
	return "option"
}

func (m *Message) name() string {
	return "message"
}

func (f *MessageField) name() string {
	return "messageField"
}

func (o *FieldOption) name() string {
	return "fieldOption"
}

func (o *OneOf) name() string {
	return "oneOf"
}

func (e *Enum) name() string {
	return "enum"
}

func (e *EnumField) name() string {
	return "enumField"
}

func (e *EnumValueOption) name() string {
	return "enumValueOption"
}

func (s *Service) name() string {
	return "service"
}

func (r *RPC) name() string {
	return "RPC"
}

func (t *ReturnType) name() string {
	return "returnType"
}
