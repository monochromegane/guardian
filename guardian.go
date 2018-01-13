package guardian

import (
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
	path     string
	handlers map[fsnotify.Op]handler
}

func newGuardian() *guardian {
	return &guardian{
		handlers: map[fsnotify.Op]handler{},
	}
}

func (g *guardian) RegisterHandler(op fsnotify.Op, hdl handler) {
	if _, ok := hdl.(*noOpCommand); ok {
		return
	}
	g.handlers[op] = hdl
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
				if handle, ok := g.handlers[event.Op]; ok {
					handle.run(event.Name)
				}
			case <-watcher.Errors:
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
