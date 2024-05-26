package assign

import (
	"fmt"
	"reflect"
	"strings"
)

// Values sets the value of each field to the supplied value.
func Values(r ValuesRequest, dst any) error {
	if err := r.Validate(); err != nil {
		return err
	}
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Pointer {
		return mustBePointerErr
	}

	dstValue = dstValue.Elem()

	for i, v := range r.NewValues {
		if v == nil {
			continue
		}
		srcValue, err := unwrapValueToValue(v)
		if err != nil {
			return err
		}

		reflectFieldName := r.FieldNames[i]
		destField, err := getReflectFieldValue(reflectFieldName, dstValue)
		if err != nil {
			return err
		}
		var assignFn AssignFunc
		if i < len(r.Assigns) {
			assignFn = r.Assigns[i]
		}
		if err = assignValue(srcValue, destField, assignFn, r.Flags); err != nil {
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

func assignValue(src, dst reflect.Value, assign AssignFunc, flags uint8) error {
	if assign != nil {
		return assign(src, dst)
	}

	switch dst.Kind() {
	case reflect.Bool:
		v, err := valueToBool(src)
		if err != nil {
			return err
		}
		dst.Set(reflect.ValueOf(v))
	case reflect.Float32:
		if kindIsFloat(src.Kind()) {
			dst.Set(reflect.ValueOf(float32(src.Float())))
		} else if flags&FuzzyFloats != 0 {
			if kindIsInt(src.Kind()) {
				dst.Set(reflect.ValueOf(float32(src.Int())))
			}
		} else {
			return fmt.Errorf("can't assign %v to %v", src.Kind(), dst.Kind())
		}
	case reflect.Float64:
		if kindIsFloat(src.Kind()) {
			dst.Set(reflect.ValueOf(src.Float()))
		} else if flags&FuzzyFloats != 0 {
			if kindIsInt(src.Kind()) {
				dst.Set(reflect.ValueOf(float64(src.Int())))
			}
		} else {
			return fmt.Errorf("can't assign %v to %v", src.Kind(), dst.Kind())
		}
	case reflect.Int:
		if kindIsInt(src.Kind()) {
			dst.Set(reflect.ValueOf(int(src.Int())))
		}
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