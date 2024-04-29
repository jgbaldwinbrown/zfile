package csvh

import (
	"io"
	"os"
	"strings"
)

func ReadLines(r io.Reader) ([]string, error) {
	b, e := io.ReadAll(r)
	if e != nil {
		return nil, e
	}
	return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n"), nil
}

func ReadLinesPath(path string) (lines []string, err error) {
	r, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer func() {
		e := r.Close()
		if err == nil {
			err = e
		}
	}()

	return ReadLines(r)
}

func ReadLinesMaybeGz(path string) (lines []string, err error) {
	r, e := OpenMaybeGz(path)
	if e != nil {
		return nil, e
	}
	defer func() {
		e := r.Close()
		if err == nil {
			err = e
		}
	}()

	return ReadLines(r)
}
