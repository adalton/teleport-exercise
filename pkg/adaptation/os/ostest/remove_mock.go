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

type RemoveRecord struct {
	Path string
}

// RemoveMock is a component that provides a mock implementation of the
// os.Remove() function.  The implementation records the paramters received
// and returns the configured NextError.
type RemoveMock struct {
	Events    []*RemoveRecord
	NextError error
}

func (w *RemoveMock) Remove(path string) error {
	w.Events = append(w.Events, &RemoveRecord{
		Path: path,
	})

	return w.NextError
}
