package pbast

import "strings"

type Syntax struct{}

func (Syntax) String() string {
	return "proto3"
}

type Import struct {
	Name       string
	Visibility Visibility
}

func NewImport(name string) *Import {
	return &Import{
		Name: name,
	}
}

func NewPublicImport(name string) *Import {
	return &Import{
		Name:       name,
		Visibility: Public,
	}
}

func NewWeakImport(name string) *Import {
	return &Import{
		Name:       name,
		Visibility: Weak,
	}
}

type Visibility int

const (
	NotSpecified Visibility = iota
	Weak
	Public
)

func (v Visibility) String() string {
	switch v {
	case Weak:
		return "weak"
	case Public:
		return "public"
	default:
		return ""
	}
}

type Package string

func NewPackage(name string) Package {
	return Package(name)
}

func NewPackageWithElements(elems []string) Package {
	lowers := make([]string, len(elems))
	for x, s := range elems {
		lowers[x] = strings.ToLower(s)
	}

	// package name is dot separated
	return Package(strings.Join(lowers, "."))
}

type Option struct {
	Name string
	// TODO: Revisit for type safety
	Value string
}

func NewOption(name, value string) *Option {
	return &Option{
		Name:  name,
		Value: value,
	}
}

type OneOf struct {
	Name   string
	Fields []*OneOfField
}

func NewOneOf(name string) *OneOf {
	return &OneOf{
		Name: name,
	}
}

func (o *OneOf) AddField(f *OneOfField) *OneOf {
	if f == nil {
		return o
	}
	o.Fields = append(o.Fields, f)
	return o
}

type OneOfField struct {
	Type    string
	Name    string
	Index   int
	Options []*Option
}

func NewOneOfField(t Type, name string, index int) *OneOfField {
	return &OneOfField{
		Type:  t.TypeName(),
		Name:  name,
		Index: index,
	}
}

func (f *OneOfField) AddOption(o *Option) *OneOfField {
	if o == nil {
		return f
	}
	f.Options = append(f.Options, o)
	return f
}
