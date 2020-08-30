package xio_test

import (
	"bytes"
	"testing"

	"github.com/ulikunitz/xio"
)

type tWriter interface {
	xio.FullWriter
	String() string
}

func testWrapper(t *testing.T, w tWriter) {
	a := []byte("123")
	n, err := w.Write(a)
	if err != nil {
		t.Fatalf("Write error %s", err)
	}
	if n != 3 {
		t.Fatalf("Write returned %d; want %d", n, 3)
	}
}

func TestWrapper(t *testing.T) {
	buf := new(bytes.Buffer)
	w := xio.WrapWriter(buf)

	const s = "foobar"

	_, err := w.WriteString(s)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}

	if g := buf.String(); g != s {
		t.Fatalf("buf returns string %q; want %q", g, s)
	}
}
