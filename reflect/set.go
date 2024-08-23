package reflect

import (
	"fmt"
	"reflect"
)

// Set sets the value of each field to the supplied value.
func Set(r SetRequest, dst any) error {
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
		var setFn SetFunc
		if i < len(r.Assigns) {
			setFn = r.Assigns[i]
		}
		if err = setValue(srcValue, destField, setFn, r.Flags); err != nil {
			return err
		}
	}
	return nil
}

func setValue(src, dst reflect.Value, assign SetFunc, flags uint8) error {
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
	case reflect.Uint64:
		if src.CanUint() {
			dst.Set(reflect.ValueOf(src.Uint()))
		} else if src.CanInt() {
			dst.Set(reflect.ValueOf(uint64(src.Int())))
		} else {
			return fmt.Errorf("field mismatch, have %v want %v", src.Kind(), dst.Kind())
		}
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
