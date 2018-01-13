package guardian

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

type testHandler struct {
	expectOp   fsnotify.Op
	expectPath string
}

func (h *testHandler) expectLog() string {
	return fsnotify.Event{Name: h.expectPath, Op: h.expectOp}.String()
}

func (h *testHandler) run(event fsnotify.Event) ([]byte, error) {
	return []byte(event.String()), nil
}

func TestGuardianRun(t *testing.T) {
	content := []byte("temporary file's content")
	dir, _ := ioutil.TempDir("", "guardian")
	defer os.RemoveAll(dir)
	tmpfn := filepath.Join(dir, "tmpfile")
	tmpfn2 := filepath.Join(dir, "tmpfile2")

	g := newGuardian()
	out := new(bytes.Buffer)
	g.logger = newLogger(out)
	g.paths = []string{dir}

	create := &testHandler{expectPath: tmpfn, expectOp: fsnotify.Create}
	g.registerHandler(fsnotify.Create, create)
	write := &testHandler{expectPath: tmpfn, expectOp: fsnotify.Write}
	g.registerHandler(fsnotify.Write, write)
	chmod := &testHandler{expectPath: tmpfn, expectOp: fsnotify.Chmod}
	g.registerHandler(fsnotify.Chmod, chmod)
	rename := &testHandler{expectPath: tmpfn, expectOp: fsnotify.Rename}
	g.registerHandler(fsnotify.Rename, rename)
	remove := &testHandler{expectPath: tmpfn2, expectOp: fsnotify.Remove}
	g.registerHandler(fsnotify.Remove, remove)

	go g.run()
	time.Sleep(500 * time.Millisecond) // Wait start monitoring
	defer g.stop()

	// Create
	ioutil.WriteFile(tmpfn, content, 0666)
	time.Sleep(100 * time.Millisecond)
	if !strings.Contains(out.String(), create.expectLog()) {
		t.Errorf("guardian output created log %s, but %s", out.String(), create.expectLog())
	}
	out.Reset()

	// Write
	ioutil.WriteFile(tmpfn, content, 0666)
	time.Sleep(100 * time.Millisecond)
	if !strings.Contains(out.String(), write.expectLog()) {
		t.Errorf("guardian output wrote log %s, but %s", out.String(), write.expectLog())
	}
	out.Reset()

	// Chmod
	os.Chmod(tmpfn, 0644)
	time.Sleep(100 * time.Millisecond)
	if !strings.Contains(out.String(), chmod.expectLog()) {
		t.Errorf("guardian output chmod log %s, but %s", out.String(), chmod.expectLog())
	}
	out.Reset()

	// Rename
	os.Rename(tmpfn, tmpfn2)
	time.Sleep(100 * time.Millisecond)
	if !strings.Contains(out.String(), rename.expectLog()) {
		t.Errorf("guardian output renamed log %s, but %s", out.String(), rename.expectLog())
	}
	out.Reset()

	// Remove
	os.Remove(tmpfn2)
	time.Sleep(100 * time.Millisecond)
	if !strings.Contains(out.String(), remove.expectLog()) {
		t.Errorf("guardian output removed log %s, but %s", out.String(), remove.expectLog())
	}
	out.Reset()
}

func TestGuardianRunRecursive(t *testing.T) {
	content := []byte("temporary file's content")
	dir, _ := ioutil.TempDir("", "guardian")
	defer os.RemoveAll(dir)

	subdir := filepath.Join(dir, "sub")
	os.Mkdir(subdir, 0777)
	tmpfn := filepath.Join(subdir, "tmpfile")

	g := newGuardian()
	out := new(bytes.Buffer)
	g.logger = newLogger(out)
	g.paths = []string{dir}

	create := &testHandler{expectPath: tmpfn, expectOp: fsnotify.Create}
	g.registerHandler(fsnotify.Create, create)

	go g.run()
	time.Sleep(500 * time.Millisecond) // Wait start monitoring
	defer g.stop()

	// Create
	ioutil.WriteFile(tmpfn, content, 0666)
	time.Sleep(100 * time.Millisecond)
	if !strings.Contains(out.String(), create.expectLog()) {
		t.Errorf("guardian output created log %s, but %s", out.String(), create.expectLog())
	}
	out.Reset()
}
