package csvh

import (
	"fmt"
	"os"
	"io"
	"compress/gzip"
	"regexp"
)

type GzReadCloser struct {
	rc io.ReadCloser
	*gzip.Reader
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

	if g.Reader, e = gzip.NewReader(g.rc); e != nil {
		g.rc.Close()
		return h(e)
	}

	return g, nil
}

func (g *GzReadCloser) Close() error {
	err := g.Reader.Close()
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
