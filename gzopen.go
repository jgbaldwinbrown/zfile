package csvh

import (
	"fmt"
	"os"
	"io"
	"compress/gzip"
	"regexp"
	"bufio"
)

type GzReadCloser struct {
	rc io.ReadCloser
	br1 *bufio.Reader
	gr *gzip.Reader
	*bufio.Reader
}

func GzOpen(path string) (*GzReadCloser, error) {
	h := func(e error) (*GzReadCloser, error) {
		return nil, fmt.Errorf("GzOpen: %w", e)
	}
	g := new(GzReadCloser)

	var e error
	if g.rc, e = os.Open(path); e != nil {
		return h(e)
	}

	g.br1 = bufio.NewReader(g.rc)

	if g.gr, e = gzip.NewReader(g.br1); e != nil {
		g.rc.Close()
		return h(e)
	}

	g.Reader = bufio.NewReader(g.gr)

	return g, nil
}

func (g *GzReadCloser) Close() error {
	err := g.gr.Close()
	if e := g.rc.Close(); err == nil {
		err = e
	}
	return err
}

var gzRe = regexp.MustCompile(`\.gz$`)

func OpenMaybeGz(path string) (io.ReadCloser, error) {
	if gzRe.MatchString(path) {
		return GzOpen(path)
	}
	return os.Open(path)
}
