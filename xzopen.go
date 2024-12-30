package zfile

import (
	"fmt"
	"os"
	"io"
	"regexp"
	"bufio"
	"github.com/ulikunitz/xz"
)

type XzReadCloser struct {
	rc io.ReadCloser
	br1 *bufio.Reader
	gr *xz.Reader
	*bufio.Reader
}

func OpenXz(path string) (*XzReadCloser, error) {
	h := func(e error) (*XzReadCloser, error) {
		return nil, fmt.Errorf("XzOpen: %w", e)
	}
	g := new(XzReadCloser)

	var e error
	if g.rc, e = os.Open(path); e != nil {
		return h(e)
	}

	g.br1 = bufio.NewReader(g.rc)

	if g.gr, e = xz.NewReader(g.br1); e != nil {
		return h(e)
	}

	g.Reader = bufio.NewReader(g.gr)

	return g, nil
}

func (g *XzReadCloser) Close() error {
	var err error
	if e := g.rc.Close(); err == nil {
		err = e
	}
	return err
}

var xzRe = regexp.MustCompile(`\.xz$`)
