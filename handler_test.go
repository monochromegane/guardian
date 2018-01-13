package guardian

import (
	"fmt"
	"testing"

	fsnotify "gopkg.in/fsnotify.v1"
)

func TestNewCommand(t *testing.T) {
	h := newCommand("")
	if _, ok := h.(*noOpCommand); !ok {
		t.Errorf("newCommand should return noOpCommand handler when empty script passed.")
	}

	h = newCommand("echo test")
	if _, ok := h.(*noOpCommand); ok {
		t.Errorf("newCommand should return command handler when script passed.")
	}
}

func TestNoOpCommandRun(t *testing.T) {
	cmd := &noOpCommand{}
	out, err := cmd.run(fsnotify.Event{Name: "dummy", Op: fsnotify.Create})
	if out != nil {
		t.Errorf("noOpCommand should return nil output, but %s", out)
	}
	if err != nil {
		t.Errorf("noOpCommand should return nil error, but %s", err)
	}
}

func TestCommandRun(t *testing.T) {
	cmd := &command{"echo -n %e %p"}

	expectPath := "dummy"
	expectEvent := fsnotify.Create
	expectOut := fmt.Sprintf("%s %s", expectEvent.String(), expectPath)

	out, err := cmd.run(fsnotify.Event{Name: expectPath, Op: expectEvent})
	if err != nil {
		t.Errorf("command should return nil error, but %s", err)
	}

	if string(out) != expectOut {
		t.Errorf("command should output %s , but %s", expectOut, string(out))
	}
}
