package ostest

import "github.com/adalton/teleport-exercise/pkg/adaptation/os"

type WriteFileRecord struct {
	Name string
	Data []byte
	Perm os.FileMode
}

type WriteFileRecorder struct {
	Events    []*WriteFileRecord
	NextError error
}

func (w *WriteFileRecorder) WriteFile(name string, data []byte, perm os.FileMode) error {
	w.Events = append(w.Events, &WriteFileRecord{
		Name: name,
		Data: data,
		Perm: perm,
	})

	return w.NextError
}
