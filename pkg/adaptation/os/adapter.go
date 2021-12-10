package os

import (
	goos "os"
)

type FileMode = goos.FileMode

// Adapter serves as a shim between between callers of standard os.* APIs
// and the functions themselves.  The default behavior is simply to dispatch
// the call to the corresponding os.* function.
//
// If the field associated with that function is non-nil, then the adapter will
// dispatch to that function instead.
type Adapter struct {
	MkdirAllFn  func(path string, perm goos.FileMode) error
	RemoveFn    func(name string) error
	WriteFileFn func(name string, data []byte, perm goos.FileMode) error
}

func (o *Adapter) MkdirAll(path string, perm goos.FileMode) error {
	fn := goos.MkdirAll

	if o != nil && o.MkdirAllFn != nil {
		fn = o.MkdirAllFn
	}

	return fn(path, perm)
}

func (o *Adapter) Remove(name string) error {
	fn := goos.Remove

	if o != nil && o.RemoveFn != nil {
		fn = o.RemoveFn
	}

	return fn(name)
}

func (o *Adapter) WriteFile(name string, data []byte, perm goos.FileMode) error {
	fn := goos.WriteFile

	if o != nil && o.WriteFileFn != nil {
		fn = o.WriteFileFn
	}

	return fn(name, data, perm)
}
