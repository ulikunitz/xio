package xio_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ulikunitz/xio"
)

type tWriter interface {
	io.Writer
	String() string
}

func testWrapper(t *testing.T, w tWriter) {
	const s = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	fw := xio.WrapWriter(w)

	i := 0
	for ; i < 3; i++ {
		if err := fw.WriteByte(s[i]); err != nil {
			t.Fatalf("w.WriteByte(s[%d]) error %s", i, err)
		}
	}
	for ; i <= 15; i += 3 {
		a := []byte(s[i : i+3])
		n, err := fw.Write(a)
		if err != nil {
			t.Fatalf("w.Write([]byte(s[%d:%d]) error %s",
				i, i+3, err)
		}
		if n != 3 {
			t.Fatalf("w.Write([]byte(s[%d:%d])"+
				" returned %d; want %d", i, i+3, n, 3)
		}
	}
	for ; i < len(s); i += 3 {
		j := i + 3
		if j > len(s) {
			j = len(s)
		}
		tmp := s[i:j]
		n, err := fw.WriteString(tmp)
		if err != nil {
			t.Fatalf("w.WriteString(s[%d:%d]) error %s",
				i, j, err)
		}
		if n != len(tmp) {
			t.Fatalf("w.WriteString(s[%d:%d]) returned %d; want %d",
				i, j, n, len(tmp))
		}
	}

	g := w.String()
	if g != s {
		t.Fatalf("w didn't return the expected string")
	}
}

type pureWriter struct {
	buf []byte
}

func (w *pureWriter) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func (w *pureWriter) String() string {
	return string(w.buf)
}

type byteWriter struct {
	pureWriter
}

func (w *byteWriter) WriteByre(c byte) error {
	w.buf = append(w.buf, c)
	return nil
}

type stringWriter struct {
	pureWriter
}

func (w *stringWriter) WriteString(s string) (n int, err error) {
	w.buf = append(w.buf, []byte(s)...)
	return len(s), nil
}

func TestWrapper(t *testing.T) {
	tests := []struct {
		s string
		w tWriter
	}{
		{"bytes.Buffer", new(bytes.Buffer)},
		{"pureWriter", new(pureWriter)},
		{"byteWriter", new(byteWriter)},
	}
	for _, tc := range tests {
		t.Run(tc.s, func(t *testing.T) {
			testWrapper(t, tc.w)
		})
	}
}
