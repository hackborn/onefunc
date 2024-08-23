package reflect

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// SetFunc handles the work of assigning the destination
// value to the source value.
type SetFunc func(src, dst reflect.Value) error

// SetJson takes a source string and a destination of any type
// and unmarshals the source to the destination.
func SetJson(src, dst reflect.Value) error {
	if src.Kind() != reflect.String {
		return fmt.Errorf("JSON format requires string source value")
	}
	if !dst.CanInterface() {
		return fmt.Errorf("JSON format requires value that CanInterface{}")
	}
	ty := reflect.TypeOf(dst.Interface())
	val := reflect.New(ty)
	err := json.Unmarshal([]byte(src.String()), val.Interface())
	dst.Set(val.Elem())
	return err
}
