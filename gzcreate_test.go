package csvh

import (
	"testing"
	"compress/gzip"
	"io"
	"bufio"
	"os"
	"fmt"
)

func Null() *os.File {
	null, e := os.Create("/dev/null")
	if e != nil {
		panic(e)
	}
	return null
}

func Writen(w io.Writer, n int) error {
	for i := 0; i < n; i++ {
		_, e := fmt.Fprintf(w, "Hello world!\n")
		if e != nil {
			return e
		}
	}
	return nil
}

func BenchmarkBuffered(b *testing.B) {
	w := Null()
	defer w.Close()
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	gw := gzip.NewWriter(bw)
	defer gw.Close()
	e := Writen(gw, b.N)
	if e != nil {
		panic(e)
	}
}

func BenchmarkUnbuffered(b *testing.B) {
	w := Null()
	defer w.Close()
	gw := gzip.NewWriter(w)
	defer gw.Close()
	e := Writen(gw, b.N)
	if e != nil {
		panic(e)
	}
}

func BenchmarkEarlyBuffered(b *testing.B) {
	w := Null()
	defer w.Close()
	gw := gzip.NewWriter(w)
	defer gw.Close()
	bw := bufio.NewWriter(gw)
	defer bw.Flush()
	e := Writen(bw, b.N)
	if e != nil {
		panic(e)
	}
}

func BenchmarkReal(b *testing.B) {
	w, e := CreateGz("/dev/null")
	if e != nil {
		panic(e)
	}
	defer w.Close()
	e = Writen(w, b.N)
	if e != nil {
		panic(e)
	}
}
