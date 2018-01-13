package guardian

import (
	"fmt"
	"os/exec"
	"strings"

	shellwords "github.com/mattn/go-shellwords"
	fsnotify "gopkg.in/fsnotify.v1"
)

type handler interface {
	run(fsnotify.Event) ([]byte, error)
}

func newCommand(script string) handler {
	if script == "" {
		return &noOpCommand{}
	}
	return &command{script: script}
}

type command struct {
	script string
}

func (cmd *command) run(event fsnotify.Event) ([]byte, error) {
	c, err := shellwords.Parse(cmd.replaceBy(event))
	if err != nil {
		return nil, err
	}
	if len(c) == 0 {
		return nil, fmt.Errorf("Script is empty.")
	}

	if len(c) == 1 {
		return exec.Command(c[0]).CombinedOutput()
	} else {
		return exec.Command(c[0], c[1:]...).CombinedOutput()
	}
}

func (cmd *command) replaceBy(event fsnotify.Event) string {
	script := strings.Replace(cmd.script, "%p", event.Name, -1)
	return strings.Replace(script, "%e", event.Op.String(), -1)
}

type noOpCommand struct {
}

func (cmd *noOpCommand) run(event fsnotify.Event) ([]byte, error) {
	return nil, nil
}
