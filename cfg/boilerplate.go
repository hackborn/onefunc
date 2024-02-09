package cfg

import (
	"reflect"
)

// Given two maps, recursively merge right into left. Adapted from
// https://stackoverflow.com/questions/22621754/how-can-i-merge-two-maps-in-go
func mergeKeys(left, right tree) tree {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			// Key exists. If my new value is a map, recurse.
			// If it's not, replace.
			if reflect.TypeOf(rightVal).Kind() == reflect.Map {
				left[key] = mergeKeys(leftVal.(tree), rightVal.(tree))
			} else {
				left[key] = rightVal
			}
		} else {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}
	return left
}
