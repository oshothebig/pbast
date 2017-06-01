package pbast

type File struct {
	Syntax   Syntax
	Package  Package
	Imports  []Import
	Options  []Option
	Messages []Message
	Enums    []Enum
	Services []Service
}

func NewFile(p Package) File {
	return File{Package: p}
}

func (f File) AddImport(i Import) File {
	nf := File(f)
	nf.Imports = append(nf.Imports, i)
	return nf
}

func (f File) AddOption(o Option) File {
	nf := File(f)
	nf.Options = append(nf.Options, o)
	return nf
}

func (f File) AddMessage(m Message) File {
	nf := File(f)
	nf.Messages = append(nf.Messages, m)
	return nf
}

func (f File) AddEnum(e Enum) File {
	nf := File(f)
	nf.Enums = append(nf.Enums, e)
	return nf
}

func (f File) AddService(s Service) File {
	nf := File(f)
	nf.Services = append(nf.Services, s)
	return nf
}
