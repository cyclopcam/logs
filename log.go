package logs

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/logging"
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

// Write logs to Google Cloud
type LogWriterGCP struct {
	Logger *logging.Logger
	Client *logging.Client
}

func (w *LogWriterGCP) Flags() LogWriterFlags {
	return 0
}

func (w *LogWriterGCP) Write(level Level, message string) {
	w.Logger.Log(logging.Entry{
		Severity: levelToGCP(level),
		Payload:  message,
	})
}

func (w *LogWriterGCP) Close() {
	w.Logger.Flush()
	w.Client.Close()
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

// The log object that you use to write logs
type Log struct {
	Output LogWriter
}

// Create a new logger
func NewLog() (*Log, error) {
	l := &Log{}
	gcpProjectID := os.Getenv("GCP_PROJECT_ID")
	gcpLogname := os.Getenv("GCP_LOGNAME")
	if gcpProjectID != "" && gcpLogname != "" {
		fmt.Printf("Logging to GCP %v / %v (you won't see further logs on stdout)\n", gcpProjectID, gcpLogname)
		client, err := logging.NewClient(context.Background(), gcpProjectID)
		if err != nil {
			return nil, fmt.Errorf("Failed to create GCP logging client: %v", err)
		}
		l.Output = &LogWriterGCP{
			Client: client,
			Logger: client.Logger(gcpLogname),
		}
	} else {
		l.Output = &LogWriterStandard{
			Output:       os.Stdout,
			EnableColors: true,
		}
		l.Infof("Logging to stdout")
	}
	return l, nil
}

// Create new log for use during unit tests
func NewTestingLog(t *testing.T) *Log {
	return &Log{
		Output: &LogWriterTest{
			T: t,
		},
	}
}

func levelToGCP(level Level) logging.Severity {
	switch level {
	case LevelDebug:
		return logging.Debug
	case LevelInfo:
		return logging.Info
	case LevelWarn:
		return logging.Warning
	case LevelError:
		return logging.Error
	case LevelCritical:
		return logging.Critical
	}
	panic("Unknown log level")
}

func levelToName(level Level) string {
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

func (l *Log) write(level Level, format string, a ...any) {
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
		prefix += levelToName(level) + " "
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

func (l *Log) Close() {
	if l.Output != nil {
		l.Output.Close()
	}
}

func (l *Log) Debugf(format string, a ...any) {
	l.write(LevelDebug, format, a...)
}

func (l *Log) Infof(format string, a ...any) {
	l.write(LevelInfo, format, a...)
}

func (l *Log) Warnf(format string, a ...any) {
	l.write(LevelWarn, format, a...)
}

func (l *Log) Errorf(format string, a ...any) {
	l.write(LevelError, format, a...)
}

func (l *Log) Criticalf(format string, a ...any) {
	l.write(LevelCritical, format, a...)
}
