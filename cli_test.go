package guardian

import (
	"bytes"
	"testing"

	fsnotify "gopkg.in/fsnotify.v1"
)

func TestCliParseArgsEmpty(t *testing.T) {
	c := &cli{err: new(bytes.Buffer)}
	_, err := c.parseArgs([]string{})
	if err == nil {
		t.Errorf("parseArg should return error, but nil")
	}
}

func TestCliParseArgs(t *testing.T) {
	c := &cli{}
	expectPath := "dir"
	g, err := c.parseArgs([]string{expectPath})
	if err != nil {
		t.Errorf("parseArg should return nil error, but %s", err)
	}
	if len(g.paths) != 1 || g.paths[0] != expectPath {
		t.Errorf("parseArg should return guardian with path %s, but %v", expectPath, g.paths)
	}
	if len(g.handlers) != 0 {
		t.Errorf("parseArg should return guardian with 0 handlers, but %d", len(g.handlers))
	}
}

func TestCliParseArgsWithOptions(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	c := &cli{out: stdout, err: stderr}
	expectPath1 := "dir1"
	expectPath2 := "dir2"
	g, err := c.parseArgs([]string{"-write", `"echo -n %e %p"`, expectPath1, expectPath2})
	if err != nil {
		t.Errorf("parseArg should return nil error, but %s", err)
	}
	if len(g.paths) != 2 {
		t.Errorf("parseArg should return guardian with path [%s, %s], but %v", expectPath1, expectPath2, g.paths)
	}
	if len(g.handlers) != 1 {
		t.Errorf("parseArg should return guardian with 1 handlers, but %d", len(g.handlers))
	}
	if h, ok := g.handlers[fsnotify.Write]; !ok {
		t.Errorf("parseArg should return guardian with write handlers.")
	} else {
		if _, ok = h.(*command); !ok {
			t.Errorf("parseArg should return guardian with command handlers.")
		}
	}
}
