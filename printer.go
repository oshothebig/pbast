package pbast

import (
	"fmt"
	"io"
	"strings"
)

type printer struct {
}

func (p *printer) Fprint(w io.Writer, n Node) {
	switch n := n.(type) {
	case File:
		p.printFile(w, n)
	case Syntax:
		p.printSyntax(w, n)
	case Import:
		p.printImport(w, n)
	case Package:
		p.printPackage(w, n)
	case Option:
		p.printOption(w, n)
	case Message:
		p.printMessage(w, n)
	case MessageField:
		p.printMessageField(w, n)
	case FieldOption:
		p.printFieldOption(w, n)
	case Enum:
		p.printEnum(w, n)
	case EnumField:
		p.printEnumField(w, n)
	case EnumValueOption:
		p.printEnumValueOption(w, n)
	case Service:
		p.printService(w, n)
	case RPC:
		p.printRPC(w, n)
	case ReturnType:
		p.printReturnType(w, n)
	}
}

func (p *printer) printFile(w io.Writer, f File) {
	// syntax
	p.Fprint(w, f.Syntax)
	// imports
	for _, i := range f.Imports {
		p.Fprint(w, i)
	}
	// packages
	for _, pkg := range f.Packages {
		p.Fprint(w, pkg)
	}
	// options
	for _, o := range f.Options {
		p.Fprint(w, o)
	}
	// messages
	for _, m := range f.Messages {
		fmt.Fprintln(w)
		p.Fprint(w, m)
	}
	// enums
	for _, e := range f.Enums {
		fmt.Fprintln(w)
		p.Fprint(w, e)
	}
	// services
	for _, s := range f.Services {
		fmt.Fprintln(w)
		p.Fprint(w, s)
	}
}

func (p *printer) printSyntax(w io.Writer, s Syntax) {
	fmt.Fprintf(w, "syntax = \"%s\";", s)
	fmt.Fprintln(w)
}

func (p *printer) printImport(w io.Writer, i Import) {
	if i.Visibility == NotSpecified {
		fmt.Fprintf(w, "import \"%s\";", i.Name)
		fmt.Fprintln(w)
		return
	}
	fmt.Fprintf(w, "import %s \"%s\";", i.Visibility, i.Name)
	fmt.Fprintln(w)
}

func (p *printer) printPackage(w io.Writer, pkg Package) {
	fmt.Fprintf(w, "package %s;", pkg.Name)
	fmt.Fprintln(w)
}

func (p *printer) printOption(w io.Writer, o Option) {
	fmt.Fprintf(w, "%s = %s;", o.Name, o.Value)
	fmt.Fprintln(w)
}

func (p *printer) printMessage(w io.Writer, m Message) {
	// name
	fmt.Fprintf(w, "message %s {", m.Name)
	fmt.Fprintln(w)

	indent := NewSpaceWriter(w, shift)
	// fields
	for _, f := range m.Fields {
		p.Fprint(indent, f)
	}
	// enums
	for _, e := range m.Enums {
		p.Fprint(indent, e)
	}
	// messages
	for _, m := range m.Messages {
		p.Fprint(indent, m)
	}

	fmt.Fprintf(w, "}")
	fmt.Fprintln(w)
}

func (p *printer) printMessageField(w io.Writer, f MessageField) {
	if f.Repeated {
		fmt.Fprintf(w, "repeated ")
	}
	fmt.Fprintf(w, "%s %s = %d", f.Type, f.Name, f.Index)

	if len(f.Options) > 0 {
		fmt.Fprint(w, " [")
		p.Fprint(w, f.Options[0])

		for _, f := range f.Options[1:] {
			fmt.Fprint(w, ", ")
			p.Fprint(w, f)
		}
		fmt.Fprint(w, "]")
	}
	fmt.Fprint(w, ";")
	fmt.Fprintln(w)
}

func (p *printer) printFieldOption(w io.Writer, o FieldOption) {
	fmt.Fprintf(w, "%s = %s", o.Name, o.Value)
}

func (p *printer) printEnum(w io.Writer, e Enum) {
	// name
	fmt.Fprintf(w, "enum %s {", e.Name)
	fmt.Fprintln(w)
	// fields
	for _, f := range e.Fields {
		p.Fprint(NewSpaceWriter(w, shift), f)
	}
	fmt.Fprintf(w, "}")
	fmt.Fprintln(w)
}

func (p *printer) printEnumField(w io.Writer, f EnumField) {
	fmt.Fprintf(w, "%s = %d", f.Name, f.Index)

	if len(f.Options) != 0 {
		fmt.Fprintf(w, " [")
		opts := []string{}
		for _, o := range f.Options {
			opts = append(opts, fmt.Sprintf("%s = %s", o.Name, o.Value))
		}
		fmt.Fprint(w, strings.Join(opts, ", "))
		fmt.Fprintf(w, "]")
	}

	fmt.Fprint(w, ";")
	fmt.Fprintln(w)
}

func (p *printer) printEnumValueOption(w io.Writer, o EnumValueOption) {
	fmt.Fprintf(w, "%s = %s", o.Name, o.Value)
}

func (p *printer) printService(w io.Writer, s Service) {
	fmt.Fprintf(w, "service %s {\n", s.Name)

	indent := NewSpaceWriter(w, shift)
	// options
	for _, o := range s.Options {
		p.Fprint(indent, o)
	}
	// RPCs
	for _, r := range s.RPCs {
		p.Fprint(indent, r)
	}

	fmt.Fprintf(w, "}")
	fmt.Fprintln(w)
}

func (p *printer) printRPC(w io.Writer, r RPC) {
	fmt.Fprintf(w, "rpc %s ", r.Name)
	p.Fprint(w, r.Input)
	fmt.Fprint(w, " returns ")
	p.Fprint(w, r.Output)
	fmt.Fprint(w, ";")
	fmt.Fprintln(w)
}

func (p *printer) printReturnType(w io.Writer, i ReturnType) {
	fmt.Fprint(w, "(")
	if i.Streamable {
		fmt.Fprint(w, "stream ")
	}
	fmt.Fprintf(w, "%s)", i.Name)
}

const shift = 2
