package pbast

import (
	"io"
	"strings"

	"bytes"

	"github.com/openconfig/goyang/pkg/indent"
)

func NewTabWriter(w io.Writer) io.Writer {
	return indent.NewWriter(w, "\t")
}

func NewSpaceWriter(w io.Writer, count int) io.Writer {
	return indent.NewWriter(w, strings.Repeat(" ", count))
}

type decorateWriter struct {
	w       io.Writer
	prefix  []byte
	postfix []byte
}

func NewDecorateWriter(w io.Writer, prefix, postfix string) io.Writer {
	return &decorateWriter{
		w:       w,
		prefix:  []byte(prefix),
		postfix: []byte(postfix),
	}
}

func (w *decorateWriter) Write(buf []byte) (int, error) {
	if _, err := w.w.Write(w.prefix); err != nil {
		return 0, err
	}

	n, err := w.w.Write(buf)
	if err != nil {
		return n, err
	}

	if _, err := w.w.Write(w.postfix); err != nil {
		return n, err
	}

	return n, nil
}

type insertWriter struct {
	w      io.Writer
	suffix []byte
	sep    []byte
}

func (w *insertWriter) Write(buf []byte) (int, error) {
	lines := bytes.SplitAfter(buf, w.sep)
	bs := bytes.Join(lines, w.suffix)
	if n, err := w.w.Write(bs); err != nil {
		return w.actualWrittenSize(n, lines), err
	}
	return len(buf), nil
}

func (w *insertWriter) actualWrittenSize(underlay int, lines [][]byte) int {
	actual := 0
	remain := underlay
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if remain <= len(line) {
			return actual + remain
		}

		if remain <= len(line)+len(w.suffix) {
			return actual + len(line)
		}

		actual += len(line)
		remain -= len(line) + len(w.suffix)
	}

	return actual
}
