package guardian

import (
	"flag"
	"fmt"
	"io"
	"runtime"

	"gopkg.in/fsnotify.v1"
)

type cli struct {
	out, err io.Writer
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

	fs := flag.NewFlagSet("guardian", flag.ContinueOnError)
	fs.SetOutput(c.err)
	fs.Usage = func() {
		fs.SetOutput(c.out)
		defer fs.SetOutput(c.err)
		fmt.Fprintf(c.out, `guardian - Monitor file changes and execute custom commands for each event.
Version: %s (%s)
Usage:
    %% guardian -create "echo %%e %%p" /path/to/monitor [...]
Options:
`, version, runtime.Version())
		fs.PrintDefaults()
	}

	var out string
	fs.StringVar(&out, "o", "", "Output file path.")

	var create, write, remove, rename, chmod string
	fs.StringVar(&create, "create", "", "handler after create operation.")
	fs.StringVar(&write, "write", "", "handler after write operation.")
	fs.StringVar(&remove, "remove", "", "handler after remove operation.")
	fs.StringVar(&rename, "rename", "", "handler after rename operation.")
	fs.StringVar(&chmod, "chmod", "", "handler after chmod operation.")

	fs.BoolVar(&g.verbose, "v", false, "Run as verbose.")
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	g.paths = fs.Args()
	if len(g.paths) == 0 {
		e := "No specify path to monitor."
		fmt.Fprintln(c.err, e)
		return nil, fmt.Errorf(e)
	}

	g.registerHandler(fsnotify.Create, newCommand(create))
	g.registerHandler(fsnotify.Write, newCommand(write))
	g.registerHandler(fsnotify.Remove, newCommand(remove))
	g.registerHandler(fsnotify.Rename, newCommand(rename))
	g.registerHandler(fsnotify.Chmod, newCommand(chmod))
	err = g.setOutput(out)
	if err != nil {
		return nil, err
	}
	return g, nil
}
