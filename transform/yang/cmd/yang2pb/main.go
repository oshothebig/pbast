package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	goyang "github.com/openconfig/goyang/pkg/yang"
	"github.com/oshothebig/pbast"
	"github.com/oshothebig/pbast/printer"
	"github.com/oshothebig/pbast/transform/yang"
)

func main() {
	var (
		fileOut  = flag.String("w", "", "output filename")
		yangPath = flag.String("path", "", "directory to add to search path")
		rewrite  = flag.Bool("r", false, "enable rewrite")
	)
	flag.Parse()

	if flag.NArg() == 0 {
		os.Exit(1)
	}

	config := config{
		out:     output(*fileOut),
		err:     os.Stderr,
		path:    []string{*yangPath},
		module:  flag.Arg(0),
		rewrite: *rewrite,
	}

	translator := translator{config: config}
	if err := translator.execute(); err != nil {
		os.Exit(1)
	}
	translator.finalize()
}

type config struct {
	out     io.Writer
	err     io.Writer
	path    []string
	module  string
	rewrite bool
}

func output(filename string) io.Writer {
	if filename == "" {
		return os.Stdout
	}

	file, err := os.Create(filename)
	if err != nil {
		// fall back to stdout
		return os.Stdout
	}
	return file
}

func (c *config) finalize() {
	if v, ok := c.out.(io.Closer); ok {
		v.Close()
	}
}

type translator struct {
	config
}

func (t *translator) execute() error {
	goyang.AddPath(t.path...)
	entry, errs := goyang.GetModule(t.module)

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(t.err, e)
		}
		return errors.New("error occurs while parsing")
	}

	protobuf := yang.Transform(entry)
	protobuf = t.rewrite(protobuf)
	printer.Fprint(t.out, protobuf)

	return nil
}

func (t *translator) rewrite(f *pbast.File) *pbast.File {
	if t.config.rewrite {
		f = yang.CompleteZeroInEnum(f)
		f = yang.AppendPrefixForEnumValueStartingWithNumber(f)
	}

	return f
}
