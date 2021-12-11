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
