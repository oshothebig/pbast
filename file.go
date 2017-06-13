package pbast

type File struct {
	Syntax   Syntax
	Package  Package
	Comment  Comment
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

func (f *File) AddType(t Type) {
	if t == nil {
		return
	}

	switch t := t.(type) {
	case *Message:
		f.AddMessage(t)
	case *Enum:
		f.AddEnum(t)
	default:
		// no-op
	}
}
