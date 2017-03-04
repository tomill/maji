package main

import "github.com/radovskyb/watcher"

func NewWatcher(dirs []string, exclude []string) (*watcher.Watcher, error) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.IgnoreHiddenFiles(true)

	for _, v := range dirs {
		if err := w.AddRecursive(v); err != nil {
			return w, err
		}
	}

	for _, v := range exclude {
		w.Ignore(v)
	}

	return w, nil
}
