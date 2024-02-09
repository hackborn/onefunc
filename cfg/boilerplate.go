package cfg

import (
	"io/fs"
)

// func readFileFS() hides the boilerplate around reading
// the full contents of a file with the new FS.
func readFileFS(fsys fs.FS, fn string) ([]byte, error) {
	f, err := fsys.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := int(fi.Size())
	b := make([]byte, size, size)
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	return b, err
}

// Given two maps, recursively merge right into left, NEVER replacing any key that already exists in left
// https://stackoverflow.com/questions/22621754/how-can-i-merge-two-maps-in-go
func mergeKeys(left, right tree) tree {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			//then we don't want to replace it - recurse
			left[key] = mergeKeys(leftVal.(tree), rightVal.(tree))
		} else {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}
	return left
}
