package ostest

import "github.com/adalton/teleport-exercise/pkg/adaptation/os"

type NilFileWriter struct{}

func (w *NilFileWriter) WriteFile(string, []byte, os.FileMode) error {
	return nil
}
