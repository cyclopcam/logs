package logs

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

// Log level
type Level int

const (
	LevelDebug    Level = iota // information that only a programmer will understand
	LevelInfo                  // information that a non-programmer might be interested in
	LevelWarn                  // speeds up tracking down issues, once you know about them
	LevelError                 // should not have happened
	LevelCritical              // wake somebody up
)

// Log writer flags
type LogWriterFlags int

const (
	LogWriterFlagNeedNewline LogWriterFlags = 1 << iota // If you should end the message with a \n
	LogWriterFlagWantDate                               // If you should embed the date at the start of the message
	LogWriterFlagWantLevel                              // If you should embed the level at the start of the message
	LogWriterFlagWantColors                             // If you should embed color escape codes
)

// LogWriter is the low level object that writes the logs
type LogWriter interface {
	Flags() LogWriterFlags
	Write(level Level, message string)
	Close()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

// Write logs to a standard file such as stdout
type LogWriterStandard struct {
	Output       io.Writer
	EnableColors bool
}

func (w *LogWriterStandard) Flags() LogWriterFlags {
	f := LogWriterFlagWantLevel | LogWriterFlagWantDate | LogWriterFlagNeedNewline
	if w.EnableColors {
		f |= LogWriterFlagWantColors
	}
	return f
}

func (w *LogWriterStandard) Write(level Level, message string) {
	w.Output.Write([]byte(message))
}

func (w *LogWriterStandard) Close() {
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

// Write logs during unit tests
type LogWriterTest struct {
	T *testing.T
}

func (w *LogWriterTest) Flags() LogWriterFlags {
	return LogWriterFlagWantLevel | LogWriterFlagWantDate
}

func (w *LogWriterTest) Write(level Level, message string) {
	w.T.Log(message)
}

func (w *LogWriterTest) Close() {
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

// We keep this interface, but in practice we only use one implementation, which is Logger.
type Log interface {
	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
	Criticalf(format string, a ...any)
	Close()
	LogWriter() LogWriter // Get the underlying LogWriter
}

// The log object that you use to write logs
type Logger struct {
	Output LogWriter
}

// Create a new logger
func NewLog() (Log, error) {
	l := &Logger{}
	l.Output = &LogWriterStandard{
		Output:       os.Stdout,
		EnableColors: true,
	}
	l.Infof("Logging to stdout")
	return l, nil
}

// Create new log for use during unit tests
func NewTestingLog(t *testing.T) Log {
	return &Logger{
		Output: &LogWriterTest{
			T: t,
		},
	}
}

func LevelToName(level Level) string {
	switch level {
	case LevelDebug:
		return "Debug"
	case LevelInfo:
		return "Info"
	case LevelWarn:
		return "Warning"
	case LevelError:
		return "Error"
	case LevelCritical:
		return "Critical"
	}
	panic("Unknown log level")
}

func (l *Logger) write(level Level, format string, a ...any) {
	flags := l.Output.Flags()
	prefix := ""
	suffix := ""
	if flags&LogWriterFlagWantDate != 0 {
		tm := time.Now().UTC().Format("2006-01-02 15:04:05.999999")
		for len(tm) < 26 {
			tm += "0"
		}
		prefix += tm + " "
	}
	if flags&LogWriterFlagWantLevel != 0 {
		prefix += LevelToName(level) + " "
	}
	if flags&LogWriterFlagWantColors != 0 && level > LevelInfo {
		prefix = "\033[0;33m" + prefix // yellow
		if level >= LevelError {
			prefix = "\033[0;31m" + prefix // red
		}
		suffix = "\033[0m"
	}
	if flags&LogWriterFlagNeedNewline != 0 {
		suffix += "\n"
	}
	l.Output.Write(level, prefix+fmt.Sprintf(format, a...)+suffix)
}

func (l *Logger) Close() {
	if l.Output != nil {
		l.Output.Close()
	}
}

func (l *Logger) LogWriter() LogWriter {
	return l.Output
}

func (l *Logger) Debugf(format string, a ...any) {
	l.write(LevelDebug, format, a...)
}

func (l *Logger) Infof(format string, a ...any) {
	l.write(LevelInfo, format, a...)
}

func (l *Logger) Warnf(format string, a ...any) {
	l.write(LevelWarn, format, a...)
}

func (l *Logger) Errorf(format string, a ...any) {
	l.write(LevelError, format, a...)
}

func (l *Logger) Criticalf(format string, a ...any) {
	l.write(LevelCritical, format, a...)
}
