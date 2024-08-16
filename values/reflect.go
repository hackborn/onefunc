package values

import (
	"fmt"
	"reflect"
	"strings"
)

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

func kindIsFloat(kind reflect.Kind) bool {
	return kind == reflect.Float32 || kind == reflect.Float64
}

func kindIsInt(kind reflect.Kind) bool {
	return kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int32 || kind == reflect.Int64
}

func unwrapValueToAny(v any) (any, error) {
	vt, err := unwrapValueToValue(v)
	if err != nil {
		return nil, err
	}
	return vt.Interface(), nil
}

func unwrapValueToValue(v any) (reflect.Value, error) {
	va := reflect.ValueOf(v)
	return doUnwrapValueToValue(va)
}

func doUnwrapValueToValue(v reflect.Value) (reflect.Value, error) {
	switch v.Kind() {
	// Directly handle the cases that need to unwrap,
	// everything else is returned as the final value.
	//	case reflect.String, reflect.Struct:
	//		return v, nil
	case reflect.Interface:
		return doUnwrapValueToValue(v.Elem())
	case reflect.Ptr:
		return doUnwrapValueToValue(v.Elem())
	default:
		return v, nil
	}
	// return reflect.Value{}, unhandledValueTypeErr
}

func valueToBool(src reflect.Value) (bool, error) {
	switch src.Kind() {
	case reflect.Bool:
		return src.Bool(), nil
	case reflect.String:
		s := strings.ToLower(src.Interface().(string))
		if s == "t" || s == "true" {
			return true, nil
		} else if s == "f" || s == "false" {
			return false, nil
		}
	}
	return false, fmt.Errorf("unsupported bool conversion on type %T", src.Interface())
}

func setField(dst reflect.Value, srcType reflect.Type, srcValue reflect.Value, i int) error {
	srcValueField := srcValue.Field(i)
	if !srcValueField.CanInterface() {
		return nil
	}
	typeField := srcType.Field(i)
	dstValueField := dst.FieldByName(typeField.Name)
	if !dstValueField.IsValid() || !dstValueField.CanInterface() {
		return nil
	}
	if srcValueField.Kind() == dstValueField.Kind() {
		dstValueField.Set(srcValueField)
	}
	return nil
}
