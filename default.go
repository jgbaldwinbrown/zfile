package csvh

import (
	"fmt"
	"errors"
)

func GetDefault[M ~map[K]V, K comparable, V any](m M, key K, def V) V {
	val, ok := m[key]
	if !ok {
		return def
	}
	return val
}

var ErrMap = errors.New("invalid map access")

func MustGet[M ~map[K]V, K comparable, V any](m M, key K) V {
	val, ok := m[key]
	if !ok {
		panic(fmt.Errorf("map does not contain key %v: %w", key, ErrMap))
	}
	return val
}

func Must(e error) {
	if e != nil {
		panic(e)
	}
}
