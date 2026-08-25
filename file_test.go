package logs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogWriterFileRollover(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "mylog.log")
	w, err := NewLogWriterFile(filename, 50)
	if err != nil {
		t.Fatal(err)
	}

	w.Write(LevelInfo, "first")
	w.Write(LevelWarn, "second")
	w.Close()

	old, err := os.ReadFile(filename + ".old")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(old), "Info first\n") {
		t.Fatalf("old log does not contain first message: %q", old)
	}
	if !strings.Contains(string(old), "Warning second\n") {
		t.Fatalf("old log does not contain rollover message: %q", old)
	}

	current, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if len(current) != 0 {
		t.Fatalf("new log is not empty: %q", current)
	}

	w, err = NewLogWriterFile(filename, 1)
	if err != nil {
		t.Fatal(err)
	}
	w.Write(LevelError, "third")
	w.Close()

	old, err = os.ReadFile(filename + ".old")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(old), "Error third\n") {
		t.Fatalf("old log was not replaced: %q", old)
	}
}
