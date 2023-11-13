package csvh

import (
	"fmt"
)

type ScanFunc func(s string, ptr any) error

func ScanF(line []string, f ScanFunc, ptrs ...any) (n int, err error) {
	if len(line) < len(ptrs) {
		return 0, fmt.Errorf("ScanF: len(line) %v < len(ptrs) %v", len(line), len(ptrs))
	}

	for i, p := range ptrs {
		e := f(line[i], p)
		if e != nil {
			return n, fmt.Errorf("ScanF: %w", e)
		}
		n++
	}

	return n, nil
}

func ScanDefault(s string, ptr any) error {
	if sp, ok := ptr.(*string); ok {
		*sp = s
		return nil
	}

	if _, e := fmt.Sscanf(s, "%v", ptr); e != nil {
		return fmt.Errorf("ScanDefault: %w", e)
	}
	return nil
}

func Scan(line []string, ptrs ...any) (n int, err error) {
	return ScanF(line, ScanDefault, ptrs...)
}

type PrintF func(v any) string

func AppendF(buf []string, f PrintF, vals ...any) []string {
	for _, val := range vals {
		buf = append(buf, f(val))
	}
	return buf
}

func PrintDefault(v any) string {
	return fmt.Sprintf("%v", v)
}

func Append(buf []string, vals ...any) []string {
	return AppendF(buf, PrintDefault, vals...)
}
