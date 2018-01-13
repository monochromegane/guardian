package guardian

import (
	"io"
	"log"
	"os"

	"github.com/farmergreg/rfsnotify"
	"gopkg.in/fsnotify.v1"
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
	paths    []string
	handlers map[fsnotify.Op]handler
	logger   *log.Logger
	verbose  bool
}

func newGuardian() *guardian {
	return &guardian{
		handlers: map[fsnotify.Op]handler{},
		logger:   newLogger(os.Stdout),
	}
}

func newLogger(out io.Writer) *log.Logger {
	return log.New(out, "[Guardian] ", log.LstdFlags)
}

func (g *guardian) registerHandler(op fsnotify.Op, hdl handler) {
	if _, ok := hdl.(*noOpCommand); ok {
		return
	}
	g.handlers[op] = hdl
}

func (g *guardian) setOutput(path string) error {
	if path == "" {
		g.logger = newLogger(os.Stdout)
		return nil
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	g.logger = newLogger(file)
	return nil
}

func (g *guardian) run() error {
	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if g.verbose {
					g.logger.Println(event)
				}
				if handle, ok := g.handlers[event.Op]; ok {
					out, err := handle.run(event)
					if err != nil {
						g.logger.Println(err)
					} else {
						g.logger.Printf("%s", out)
					}
				}
			case err := <-watcher.Errors:
				g.logger.Println(err)
			}
		}
	}()

	for _, path := range g.paths {
		err = watcher.AddRecursive(path)
		if err != nil {
			return err
		}
	}
	<-done
	return nil
}
