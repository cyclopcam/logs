package logs

import (
	"bytes"
	"testing"
)

func TestPrefix(t *testing.T) {
	var buf bytes.Buffer
	log := &Logger{
		Output: &LogWriterStandard{
			Output:       &buf,
			EnableColors: false,
			EnableDate:   false, // easier to test without date
			EnableLevel:  true,
		},
		Prefix: "",
	}
	prefixedLog := NewPrefixLogger(log, "[myprefix]")

	prefixedLog.Infof("Hello %v", "world")
	expected := "Info [myprefix] Hello world\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}
