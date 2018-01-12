package guardian

import (
	"flag"
	"io"
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
	g := &guardian{}
	flag.StringVar(&g.path, "p", "", "path")
	flag.Parse()
	return g, nil
}
