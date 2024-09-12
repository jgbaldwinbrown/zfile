package csvh

import (
	"testing"
	"bufio"

	"github.com/ulikunitz/xz"
)

func BenchmarkBufferedXz(b *testing.B) {
	w := Null()
	defer w.Close()
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	gw, e := xz.NewWriter(bw)
	if e != nil {
		panic(e)
	}
	defer gw.Close()
	e = Writen(gw, b.N)
	if e != nil {
		panic(e)
	}
}

func BenchmarkUnbufferedXz(b *testing.B) {
	w := Null()
	defer w.Close()
	gw, e := xz.NewWriter(w)
	if e != nil {
		panic(e)
	}
	defer gw.Close()
	e = Writen(gw, b.N)
	if e != nil {
		panic(e)
	}
}

func BenchmarkEarlyBufferedXz(b *testing.B) {
	w := Null()
	defer w.Close()
	gw, e := xz.NewWriter(w)
	if e != nil {
		panic(e)
	}
	defer gw.Close()
	bw := bufio.NewWriter(gw)
	defer bw.Flush()
	e = Writen(bw, b.N)
	if e != nil {
		panic(e)
	}
}

func BenchmarkRealXz(b *testing.B) {
	w, e := XzCreate("/dev/null")
	if e != nil {
		panic(e)
	}
	defer w.Close()
	e = Writen(w, b.N)
	if e != nil {
		panic(e)
	}
}
