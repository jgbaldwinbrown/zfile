package zfile

import (
	"bufio"
	"fmt"
	"os"
	"io"

	"github.com/biogo/hts/bgzf"
)

type BgzWriteCloser struct {
	wc io.WriteCloser
	bw1 *bufio.Writer
	gw *bgzf.Writer
	*bufio.Writer
}

func CreateBgz(path string, threads int) (*BgzWriteCloser, error) {
	h := func(e error) (*BgzWriteCloser, error) {
		return nil, fmt.Errorf("BgzCreate: %w", e)
	}
	g := new(BgzWriteCloser)

	var e error
	if g.wc, e = os.Create(path); e != nil {
		return h(e)
	}

	g.bw1 = bufio.NewWriter(g.wc)
	g.gw = bgzf.NewWriter(g.bw1, threads)
	g.Writer = bufio.NewWriter(g.gw)

	return g, nil
}

func (g *BgzWriteCloser) Close() error {
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
