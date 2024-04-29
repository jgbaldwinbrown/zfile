package csvh

import (
	"strings"
)

func AppendFields(out []string, in, sep string) []string {
	rem := in
	var t string
	var found bool
	for t, rem, found = strings.Cut(rem, sep); found; t, rem, found = strings.Cut(rem, sep) {
		out = append(out, t)
	}
	out = append(out, rem)
	return out
}
