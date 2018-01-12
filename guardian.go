package guardian

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

const (
	exitCodeOK = iota
	exitCodeErr
)

func Run(args []string) int {
	err := (&cli{out: os.Stdout, err: os.Stderr}).run(args)
	if err != nil {
		return exitCodeErr
	}
	return exitCodeOK
}

type guardian struct {
	path string
}

func (g *guardian) run() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println("event:", event)
			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(g.path)
	if err != nil {
		return err
	}
	<-done
	return nil
}
