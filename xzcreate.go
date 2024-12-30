package csvh

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"github.com/ulikunitz/xz"
)

type XzWriteCloser struct {
	wc io.WriteCloser
	bw1 *bufio.Writer
	gw *xz.Writer
	*bufio.Writer
}

func CreateXz(path string) (*XzWriteCloser, error) {
	h := func(e error) (*XzWriteCloser, error) {
		return nil, fmt.Errorf("XzCreate: %w", e)
	}
	g := new(XzWriteCloser)

	var e error
	if g.wc, e = os.Create(path); e != nil {
		return h(e)
	}

	g.bw1 = bufio.NewWriter(g.wc)
	g.gw, e = xz.NewWriter(g.bw1)
	if e != nil {
		g.wc.Close()
		g.bw1.Flush()
		return h(e)
	}
	g.Writer = bufio.NewWriter(g.gw)

	return g, nil
}

func (g *XzWriteCloser) Close() error {
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
