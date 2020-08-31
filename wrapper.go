// Package xio extends a Writer to support WriteByte and WriteString.
//
// Sometimes you might get a Writer for which you need a WriteByte or
// WriteString method. The WrapWriter function returns a full writer
// supporting both methods. Alternatively you could use bufio.Writer,
// which rquires however to call the Flush method. The writer returned
// from WrapWriter executes all Write commands directly and doesn't
// require a Flush method.
package xio

import (
	"errors"
	"io"
)

// A full writer supports the Write, WriteByte and WriteString methods.
type FullWriter interface {
	io.Writer
	io.StringWriter
	io.ByteWriter
}

// wrapper converts an io.Writer to an io.FullWriter.
type wrapper struct {
	io.Writer
	bw        io.ByteWriter
	sw        io.StringWriter
	byteBuf   []byte
	stringBuf []byte
}

// WrapWriter converts a writer into a writer that supports Write,
// WriteByte and WriteString.
func WrapWriter(w io.Writer) FullWriter {
	if fw, ok := w.(FullWriter); ok {
		return fw
	}

	fw := &wrapper{Writer: w}

	// If the writer is a byte writer call the function.
	if bw, ok := w.(io.ByteWriter); ok {
		fw.bw = bw
	} else {
		fw.byteBuf = make([]byte, 1)
	}

	// If the writer is a string writer call the function directly.
	if sw, ok := w.(io.StringWriter); ok {
		fw.sw = sw
	} else {
		fw.stringBuf = make([]byte, 0, 32)
	}

	return fw
}

func (w *wrapper) WriteByte(c byte) error {
	if w.bw != nil {
		return w.bw.WriteByte(c)
	}
	w.byteBuf[0] = c
	n, err := w.Write(w.byteBuf)
	if n == 1 {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New("WriteByte: Write returned no error")

}

func (w *wrapper) WriteString(s string) (n int, err error) {
	if w.sw != nil {
		return w.sw.WriteString(s)
	}
	w.stringBuf = append(w.stringBuf[:0], s...)
	return w.Write(w.stringBuf)
}
