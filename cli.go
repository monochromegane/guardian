package guardian

import (
	"flag"
	"io"

	"gopkg.in/fsnotify.v1"
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

	var out string
	flag.StringVar(&out, "o", "", "Output file path.")

	var create, write, remove, rename, chmod string
	flag.StringVar(&create, "create", "", "handler after create operation.")
	flag.StringVar(&write, "write", "", "handler after write operation.")
	flag.StringVar(&remove, "remove", "", "handler after remove operation.")
	flag.StringVar(&rename, "rename", "", "handler after rename operation.")
	flag.StringVar(&chmod, "chmod", "", "handler after chmod operation.")

	flag.BoolVar(&g.verbose, "v", false, "Run as verbose.")
	flag.Parse()

	g.registerHandler(fsnotify.Create, newCommand(create))
	g.registerHandler(fsnotify.Write, newCommand(write))
	g.registerHandler(fsnotify.Remove, newCommand(remove))
	g.registerHandler(fsnotify.Rename, newCommand(rename))
	g.registerHandler(fsnotify.Chmod, newCommand(chmod))
	err := g.setOutput(out)
	if err != nil {
		return nil, err
	}
	return g, nil
}
