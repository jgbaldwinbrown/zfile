package csvh

import (
	"fmt"
	"os"
	"io"
	"compress/gzip"
)

type GzWriteCloser struct {
	wc io.WriteCloser
	*gzip.Writer
}

func GzCreate(path string) (*GzWriteCloser, error) {
	h := func(e error) (*GzWriteCloser, error) {
		return nil, fmt.Errorf("GzCreate: %w", e)
	}
	g := new(GzWriteCloser)
	var e error

	if g.wc, e = os.Create(path); e != nil {
		return h(e)
	}

	g.Writer = gzip.NewWriter(g.wc)

	return g, nil
}

func (g *GzWriteCloser) Close() error {
	err := g.Writer.Close()
	if e := g.wc.Close(); err == nil {
		err = e
	}
	return err
}

func CreateMaybeGz(path string) (io.WriteCloser, error) {
	if gzRe.MatchString(path) {
		return GzCreate(path)
	}
	return os.Create(path)
}
