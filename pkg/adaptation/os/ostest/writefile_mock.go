/*
Copyright 2021 Andy Dalton
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ostest

import "github.com/adalton/teleport-exercise/pkg/adaptation/os"

type WriteFileRecord struct {
	Name string
	Data []byte
	Perm os.FileMode
}

// WriteFileMock is a component that provides a mock implementation of the
// os.WriteFile() function.  The implementation records the paramters received
// and returns the configured NextError.
type WriteFileMock struct {
	Events    []*WriteFileRecord
	NextError error
}

func (w *WriteFileMock) WriteFile(name string, data []byte, perm os.FileMode) error {
	w.Events = append(w.Events, &WriteFileRecord{
		Name: name,
		Data: data,
		Perm: perm,
	})

	return w.NextError
}
