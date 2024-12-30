package csvh

import (
	"io"
	"bufio"
	"fmt"
	"os"
	"compress/bzip2"
)

type Bzip2ReadCloser struct {
	rc io.ReadCloser
	br1 *bufio.Reader
	gr io.Reader
	*bufio.Reader
}

func OpenBzip2(path string) (*Bzip2ReadCloser, error) {
	h := func(e error) (*Bzip2ReadCloser, error) {
		return nil, fmt.Errorf("Bzip2Open: %w", e)
	}
	g := new(Bzip2ReadCloser)

	var e error
	if g.rc, e = os.Open(path); e != nil {
		return h(e)
	}

	g.br1 = bufio.NewReader(g.rc)
	g.gr = bzip2.NewReader(g.br1)
	g.Reader = bufio.NewReader(g.gr)

	return g, nil
}

func (g *Bzip2ReadCloser) Close() error {
	var err error
	if e := g.rc.Close(); err == nil {
		err = e
	}
	return err
}
