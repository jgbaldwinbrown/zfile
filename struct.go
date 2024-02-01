package csvh

import (
	"reflect"
)

func Fields(strp any) []reflect.Value {
	sv := reflect.ValueOf(strp).Elem()
	n := sv.NumField()
	fields := make([]reflect.Value, 0, n)
	for i := 0; i < n; i++ {
		fields = append(fields, sv.Field(i))
	}
	return fields
}

func AppendVals(dst []any, fields ...reflect.Value) []any {
	for _, f := range fields {
		dst = append(dst, f.Interface())
	}
	return dst
}

func Ptrs(strp any) []any {
	sv := reflect.ValueOf(strp).Elem()
	n := sv.NumField()
	ptrs := make([]any, 0, n)
	for i := 0; i < n; i++ {
		ptrs = append(ptrs, sv.Field(i).Addr().Interface())
	}
	return ptrs
}
