package zfile

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"compress/gzip"
)

type GzWriteCloser struct {
	wc io.WriteCloser
	bw1 *bufio.Writer
	gw *gzip.Writer
	*bufio.Writer
}

func CreateGz(path string) (*GzWriteCloser, error) {
	h := func(e error) (*GzWriteCloser, error) {
		return nil, fmt.Errorf("GzCreate: %w", e)
	}
	g := new(GzWriteCloser)

	var e error
	if g.wc, e = os.Create(path); e != nil {
		return h(e)
	}

	g.bw1 = bufio.NewWriter(g.wc)
	g.gw = gzip.NewWriter(g.bw1)
	g.Writer = bufio.NewWriter(g.gw)

	return g, nil
}

func (g *GzWriteCloser) Close() error {
	err := g.Writer.Flush()
	if e := g.gw.Close(); err == nil {
		err = e
	}
	if e := g.bw1.Flush(); err == nil {
		err = e
	}
	if e := g.wc.Close(); err == nil {
		err = e
	}
	return err
}

func Create(path string) (io.WriteCloser, error) {
	if gzRe.MatchString(path) {
		return CreateGz(path)
	}
	if xzRe.MatchString(path) {
		return CreateXz(path)
	}
	return CreateUncompressed(path)
}

func SupportedCreate() []string {
	return []string{".gz", ".xz"}
}
