package pbast

type Enum struct {
	Name   string
	Fields []EnumField
}

func NewEnum(name string) Enum {
	return Enum{
		Name: name,
	}
}

func (e Enum) AddField(f EnumField) Enum {
	ne := Enum(e)
	ne.Fields = append(ne.Fields, f)
	return ne
}

type EnumField struct {
	Name    string
	Index   int
	Options []EnumValueOption
}

type EnumValueOption struct {
	Name  string
	Value string
}

func NewEnumField(name string, index int) EnumField {
	return EnumField{
		Name:  name,
		Index: index,
	}
}

func NewEnumValueOption(name, value string) EnumValueOption {
	return EnumValueOption{
		Name:  name,
		Value: value,
	}
}

func (f EnumField) AddOption(o EnumValueOption) EnumField {
	nf := EnumField(f)
	nf.Options = append(nf.Options, o)
	return nf
}
