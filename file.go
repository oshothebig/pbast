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
	f.Imports = append(f.Imports, i)
	return f
}

func (f *File) AddOption(o *Option) *File {
	f.Options = append(f.Options, o)
	return f
}

func (f *File) AddMessage(m *Message) *File {
	f.Messages = append(f.Messages, m)
	return f
}

func (f *File) AddEnum(e *Enum) *File {
	f.Enums = append(f.Enums, e)
	return f
}

func (f *File) AddService(s *Service) *File {
	f.Services = append(f.Services, s)
	return f
}
