package csvh

func GetDefault[M ~map[K]V, K comparable, V any](m M, key K, def V) V {
	val, ok := m[key]
	if !ok {
		return def
	}
	return val
}
