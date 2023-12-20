package assign

import (
	"fmt"
	"reflect"
)

type ValuesRequest struct {
	FieldNames []string
	NewValues  []any
}

func (r ValuesRequest) Validate() error {
	if len(r.FieldNames) != len(r.NewValues) {
		return fmt.Errorf("Size mismatch (%v names but %v values)", len(r.FieldNames), len(r.NewValues))
	}
	return nil
}

// Values sets the value of each field to the supplied value.
// func Values(r ValuesRequest, tox *_toxMetadata, tags, scanned []any, dst any) error {
func Values(r ValuesRequest, dst any) error {
	if err := r.Validate(); err != nil {
		return err
	}
	reflectValue := reflect.ValueOf(dst)
	if reflectValue.Kind() != reflect.Pointer {
		return mustBePointerErr
	}
	reflectValue = reflectValue.Elem()

	for i, v := range r.NewValues {
		if v == nil {
			continue
		}
		value, err := unwrapValue(v)
		if err != nil {
			return err
		}

		reflectFieldName := r.FieldNames[i]
		destField, err := getReflectFieldValue(reflectFieldName, reflectValue)
		if err != nil {
			return err
		}

		// Technically we shouldn't fully unwrap the value,
		// so I don't need to get it again here.
		err = assignValue(reflect.ValueOf(value), destField)
		if err != nil {
			return err
		}
	}
	return nil
}

func getReflectFieldValue(fieldName string, structValue reflect.Value) (reflect.Value, error) {
	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return reflect.Value{}, fmt.Errorf("no field for %v", fieldName)
	}
	if !field.CanSet() {
		return reflect.Value{}, fmt.Errorf("can't set field for %v", fieldName)
	}
	return field, nil
}

func assignValue(src, dst reflect.Value) error {
	switch dst.Kind() {
	case reflect.Int:
		dst.Set(reflect.ValueOf(int(src.Int())))
	case reflect.Int8:
		dst.Set(reflect.ValueOf(int8(src.Int())))
	case reflect.Int32:
		dst.Set(reflect.ValueOf(int32(src.Int())))
	case reflect.Int64:
		dst.Set(reflect.ValueOf(src.Int()))
	case reflect.String:
		if src.Kind() != reflect.String {
			return fmt.Errorf("field mismatch, have %v want %v", src.Kind(), dst.Kind())
		}
		dst.Set(reflect.ValueOf(src.String()))
	default:
		return fmt.Errorf("unsupported field type %v", dst.Kind())
	}
	return nil
}

func unwrapValue(v any) (any, error) {
	va := reflect.ValueOf(v)
	return unwrapReflectValue(va)
}

func unwrapReflectValue(v reflect.Value) (any, error) {
	switch v.Kind() {
	case reflect.String, reflect.Struct:
		return v.Interface(), nil
	case reflect.Interface:
		return unwrapReflectValue(v.Elem())
	case reflect.Ptr:
		return unwrapReflectValue(v.Elem())
	}
	return nil, unhandledValueTypeErr
}
