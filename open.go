package zfile

import (
	"os"
	"bufio"
	"io"
)

type ReadCloser struct {
	r io.ReadCloser
	b *bufio.Reader
}

func OpenUncompressed(path string) (*ReadCloser, error) {
	r := new(ReadCloser)
	var e error
	r.r, e = os.Open(path)
	if e != nil {
		return nil, e
	}
	r.b = bufio.NewReader(r.r)
	return r, nil
}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	return r.b.Read(p)
}

func (r *ReadCloser) Close() error {
	return r.r.Close()
}

type WriteCloser struct {
	w io.WriteCloser
	b *bufio.Writer
}

func CreateUncompressed(path string) (*WriteCloser, error) {
	r := new(WriteCloser)
	var e error
	r.w, e = os.Create(path)
	if e != nil {
		return nil, e
	}
	r.b = bufio.NewWriter(r.w)
	return r, nil
}

func (r *WriteCloser) Write(p []byte) (n int, err error) {
	return r.b.Write(p)
}

func (r *WriteCloser) Close() error {
	e1 := r.b.Flush()
	e2 := r.w.Close()
	if e1 != nil {
		return e1
	}
	return e2
}
