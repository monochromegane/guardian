package guardian

import (
	"flag"
	"io"

	"github.com/fsnotify/fsnotify"
)

type cli struct {
	out io.Writer
	err io.Writer
}

func (c *cli) run(args []string) error {
	g, err := c.parseArgs(args)
	if err != nil {
		return err
	}
	return g.run()
}

func (c *cli) parseArgs(args []string) (*guardian, error) {
	g := newGuardian()
	flag.StringVar(&g.path, "p", "", "path")

	var create, write, remove, rename, chmod string
	flag.StringVar(&create, "create", "", "handler after create operation.")
	flag.StringVar(&write, "write", "", "handler after write operation.")
	flag.StringVar(&remove, "remove", "", "handler after remove operation.")
	flag.StringVar(&rename, "rename", "", "handler after rename operation.")
	flag.StringVar(&chmod, "chmod", "", "handler after chmod operation.")
	flag.Parse()

	g.RegisterHandler(fsnotify.Create, newCommand(create))
	g.RegisterHandler(fsnotify.Write, newCommand(write))
	g.RegisterHandler(fsnotify.Remove, newCommand(remove))
	g.RegisterHandler(fsnotify.Rename, newCommand(rename))
	g.RegisterHandler(fsnotify.Chmod, newCommand(chmod))
	return g, nil
}
