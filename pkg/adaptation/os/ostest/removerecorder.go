package ostest

type RemoveRecord struct {
	Path string
}

type RemoveRecorder struct {
	Events    []*RemoveRecord
	NextError error
}

func (w *RemoveRecorder) Remove(path string) error {
	w.Events = append(w.Events, &RemoveRecord{
		Path: path,
	})

	return w.NextError
}
