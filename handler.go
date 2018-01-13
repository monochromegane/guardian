package guardian

import (
	"fmt"
	"os/exec"
	"strings"

	shellwords "github.com/mattn/go-shellwords"
)

type handler interface {
	run(string) ([]byte, error)
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

func (cmd *command) run(path string) ([]byte, error) {
	c, err := shellwords.Parse(strings.Replace(cmd.script, "%p", path, -1))
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

type noOpCommand struct {
}

func (cmd *noOpCommand) run(path string) ([]byte, error) {
	return nil, nil
}
