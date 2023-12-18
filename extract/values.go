package extract

import (
	"reflect"
)

// Slice iterates the fields in a struct, adding the
// results to a slice based on a handler.
func Slice(s any, h Handler, connector, assignment string) []any {
	sh := &sliceHandler{connector: connector,
		assignment: assignment}
	t := tail(h)
	if t != nil {
		t.SetNext(sh)
	} else {
		h = t
	}
	Values(s, h)
	return sh.result
}

// Values iterates the fields in a struct, sending the results
// to a handler.
func Values(s any, h Handler) {
	s = getStruct(s)
	rType := reflect.TypeOf(s)
	rValue := reflect.ValueOf(s)

	if rType.Kind() == reflect.Struct {
		len := rType.NumField()

		for i := 0; i < len; i++ {
			typeField := rType.Field(i)
			valueField := rValue.Field(i)

			h.Handle(typeField.Name, valueField.Interface())
		}
	}
}

// getStruct answers s as a struct via reflection. If s is
// a pointer to a struct, that gets unwrapped, and the struct is returned.
func getStruct(s any) any {
	rType := reflect.TypeOf(s)
	if rType.Kind() != reflect.Ptr {
		return s
	}
	v := reflect.Indirect(reflect.ValueOf(s))
	elem := v.Interface()
	return elem
}
