package reflect

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
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
			return setNumberValue[float32](src, dst)
		} else {
			return fmt.Errorf("can't assign %v to %v", src.Kind(), dst.Kind())
		}
	case reflect.Float64:
		if kindIsFloat(src.Kind()) {
			dst.Set(reflect.ValueOf(src.Float()))
		} else if flags&FuzzyFloats != 0 {
			return setNumberValue[float64](src, dst)
		} else {
			return fmt.Errorf("can't assign %v to %v", src.Kind(), dst.Kind())
		}
	case reflect.Int:
		return setNumberValue[int](src, dst)
	case reflect.Int8:
		return setNumberValue[int8](src, dst)
	case reflect.Int32:
		return setNumberValue[int32](src, dst)
	case reflect.Int64:
		return setNumberValue[int64](src, dst)
	case reflect.Uint64:
		return setNumberValue[uint64](src, dst)
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

func setNumberValue[T constraints.Integer | constraints.Float](src, dst reflect.Value) error {
	if src.CanUint() {
		dst.Set(reflect.ValueOf(T(src.Uint())))
	} else if src.CanInt() {
		dst.Set(reflect.ValueOf(T(src.Int())))
	} else if src.CanFloat() {
		dst.Set(reflect.ValueOf(T(src.Float())))
	} else {
		return fmt.Errorf("field mismatch, have %v want %v", src.Kind(), dst.Kind())
	}
	return nil
}
