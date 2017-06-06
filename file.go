package pbast

type File struct {
	Syntax   Syntax
	Package  Package
	Imports  []*Import
	Options  []*Option
	Messages []*Message
	Enums    []*Enum
	Services []*Service
}

func NewFile(p Package) *File {
	return &File{Package: p}
}

func (f *File) AddImport(i *Import) *File {
	if i == nil {
		return f
	}
	f.Imports = append(f.Imports, i)
	return f
}

func (f *File) AddOption(o *Option) *File {
	if o == nil {
		return f
	}
	f.Options = append(f.Options, o)
	return f
}

func (f *File) AddMessage(m *Message) *File {
	if m == nil {
		return f
	}
	f.Messages = append(f.Messages, m)
	return f
}

func (f *File) AddEnum(e *Enum) *File {
	if e == nil {
		return f
	}
	f.Enums = append(f.Enums, e)
	return f
}

func (f *File) AddService(s *Service) *File {
	if s == nil {
		return f
	}
	f.Services = append(f.Services, s)
	return f
}

func (f *File) AddType(t Type) *File {
	if t == nil {
		return f
	}

	switch t := t.(type) {
	case *Message:
		return f.AddMessage(t)
	case *Enum:
		return f.AddEnum(t)
	default:
		return f
	}
}
