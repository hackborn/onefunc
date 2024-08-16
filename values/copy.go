package values

import (
	"cmp"
	"reflect"
)

// Copy sets all fields in dst to any fields in src
// with the same name. It is a convenience and small
// optimization over using Set().
// WIP: Currently just a simple struct-to-struct
// top-level copy. Additional details will be added
// as needed.
func Copy(dst, src any) error {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Pointer {
		return mustBePointerErr
	}
	dstValue = dstValue.Elem()
	srcValue := getStruct(src)
	srcRType := reflect.TypeOf(srcValue)
	srcRValue := reflect.ValueOf(srcValue)

	var err error
	if srcRType.Kind() == reflect.Struct {
		len := srcRType.NumField()

		for i := 0; i < len; i++ {
			err = cmp.Or(err, setField(dstValue, srcRType, srcRValue, i))
		}
	}

	return err
}
