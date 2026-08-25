package logs

import (
	"bytes"
	"testing"
)

func TestPrefix(t *testing.T) {
	var buf bytes.Buffer
	log := &Logger{
		Output: &LogWriterStandard{
			Output: &buf,
			Flags:  LogWriterFlagWantLevel | LogWriterFlagNeedNewline,
		},
	}
	prefixedLog := NewPrefixLogger(log, "[myprefix]")

	prefixedLog.Infof("Hello %v", "world")
	expected := "Info [myprefix] Hello world\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestTee(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	log1 := &Logger{
		Output: &LogWriterStandard{
			Output: &buf1,
			Flags:  LogWriterFlagWantLevel | LogWriterFlagNeedNewline,
		},
	}
	log2 := &Logger{
		Output: &LogWriterStandard{
			Output: &buf2,
			Flags:  LogWriterFlagWantLevel | LogWriterFlagNeedNewline,
		},
	}

	teeLog := Tee(log1, log2)
	teeLog.Infof("Hello %v", "world")
	expected := "Info Hello world\n"
	if buf1.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf1.String())
	}
	if buf2.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf2.String())
	}
}
