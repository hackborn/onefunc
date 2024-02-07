package extract

import (
	"reflect"
)

// AsMap iterates the fields in a struct, adding the
// results to a map based on a handler and returning the map.
// The handler can be a chain. The final element of the
// chain must be a Mapper, but you don't need to
// include it: If the chain doesn't end in a Mapper, then
// one will be added and provided the MapOpts.
func AsMap(s any, h Handler, opts *MapOpts) map[string]any {
	mapper, ok := getLast[Mapper](h)
	if !ok {
		if opts == nil {
			opts = &MapOpts{}
		}
		h = NewChain(h, opts)
		if mapper, ok = getLast[Mapper](h); !ok {
			return nil
		}
	}
	From(s, h)
	return mapper.Map()
}

// AsSlice iterates the fields in a struct, adding the
// results to a slice based on a handler and returning the slice.
// The handler can be a chain. The final element of the
// chain must be a Slicer, but you don't need to
// include it: If the chain doesn't end in a Slicer, then
// one will be added and provided the SliceOpts.
func AsSlice(s any, h Handler, opts *SliceOpts) []any {
	slicer, ok := getLast[Slicer](h)
	if !ok {
		if opts == nil {
			opts = &SliceOpts{}
		}
		h = NewChain(h, opts)
		if slicer, ok = getLast[Slicer](h); !ok {
			return nil
		}
	}
	From(s, h)
	return slicer.Slice()
}

// From iterates the fields in a struct, sending the results
// to a handler.
func From(s any, h Handler) {
	s = getStruct(s)
	rType := reflect.TypeOf(s)
	rValue := reflect.ValueOf(s)

	if rType.Kind() == reflect.Struct {
		len := rType.NumField()

		for i := 0; i < len; i++ {
			valueField := rValue.Field(i)
			if valueField.CanInterface() {
				typeField := rType.Field(i)
				h.Handle(typeField.Name, valueField.Interface())
			}
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
