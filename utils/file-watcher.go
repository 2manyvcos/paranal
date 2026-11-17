package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	watcher  *fsnotify.Watcher
	events   chan fsnotify.Event
	errors   chan error
	done     chan struct{}
	doneOnce sync.Once
}

func (w *FileWatcher) Events() <-chan fsnotify.Event {
	return w.events
}

func (w *FileWatcher) Errors() <-chan error {
	return w.errors
}

func (w *FileWatcher) Shutdown() {
	w.doneOnce.Do(func() {
		close(w.done)
		close(w.events)
		close(w.errors)
	})

	w.watcher.Close()
}

func (w *FileWatcher) add(path string) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return w.watcher.Add(path)
		}

		return nil
	})
}

func (w *FileWatcher) handle() {
	for {
		select {
		case <-w.done:
			return

		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					err := w.add(event.Name)
					if err != nil {
						w.errors <- err
					}
				}
			}

			select {
			case w.events <- event:
			case <-w.done:
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

			select {
			case w.errors <- err:
			case <-w.done:
			}
		}
	}
}

func NewFileWatcher(path string) (*FileWatcher, error) {
	var w FileWatcher
	w.events = make(chan fsnotify.Event)
	w.errors = make(chan error)
	w.done = make(chan struct{})

	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go w.handle()

	err = w.add(path)
	if err != nil {
		w.Shutdown()
		return nil, err
	}

	return &w, nil
}
