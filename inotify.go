package main

import (
	"github.com/fsnotify/fsnotify"
)

func watchDirs(paths []string, cb func(string)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, path := range paths {
		watcher.Add(path)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					cb(event.Name)
				}
			case <-watcher.Errors:
				return
			}
		}
	}()

	return nil
}
