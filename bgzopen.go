package zfile

import (
	"fmt"
	"os"
	"io"
	"bufio"

	"github.com/biogo/hts/bgzf"
)

type BgzReadCloser struct {
	rc io.ReadCloser
	br1 *bufio.Reader
	gr *bgzf.Reader
	*bufio.Reader
}

func OpenBgz(path string, threads int) (*BgzReadCloser, error) {
	h := func(e error) (*BgzReadCloser, error) {
		return nil, fmt.Errorf("BgzOpen: %w", e)
	}
	g := new(BgzReadCloser)

	var e error
	if g.rc, e = os.Open(path); e != nil {
		return h(e)
	}

	g.br1 = bufio.NewReader(g.rc)

	if g.gr, e = bgzf.NewReader(g.br1, threads); e != nil {
		g.rc.Close()
		return h(e)
	}

	g.Reader = bufio.NewReader(g.gr)

	return g, nil
}

func (g *BgzReadCloser) Close() error {
	err := g.gr.Close()
	if e := g.rc.Close(); err == nil {
		err = e
	}
	return err
}
