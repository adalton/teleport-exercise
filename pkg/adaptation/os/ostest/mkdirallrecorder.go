package ostest

import "github.com/adalton/teleport-exercise/pkg/adaptation/os"

type MkdirAllRecord struct {
	Path string
	Perm os.FileMode
}

type MkdirAllRecorder struct {
	Events    []*MkdirAllRecord
	NextError error
}

func (w *MkdirAllRecorder) MkdirAll(path string, perm os.FileMode) error {
	w.Events = append(w.Events, &MkdirAllRecord{
		Path: path,
		Perm: perm,
	})

	return w.NextError
}
