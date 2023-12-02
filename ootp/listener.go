package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type FileListener struct {
	FileName string
	watcher  *fsnotify.Watcher
	messages chan string
}

func NewFileListener(fileName string, messages chan string) (*FileListener, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &FileListener{FileName: fileName, watcher: watcher, messages: messages}, nil
}

func (fl *FileListener) Watch() {
	defer fl.watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-fl.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fl.messages <- "File updated: " + event.Name
				}
			case err, ok := <-fl.watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	// Add the file to be watched.
	err := fl.watcher.Add(fl.FileName)
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}
