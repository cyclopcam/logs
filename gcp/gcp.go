package logsgcp

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
	"github.com/cyclopcam/logs/v3"
)

// Write logs to Google Cloud
type LogWriterGCP struct {
	Logger *logging.Logger
	Client *logging.Client
}

// Create a new logger that writes to GCP, if the GCP_PROJECT_ID and GCP_LOGNAME are set.
// If those environment variables are not set, then we fall back to stdout.
func NewLog() (logs.Log, error) {
	l := &logs.Logger{}
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
		return l, nil
	} else {
		return logs.NewLog()
	}
}

func (w *LogWriterGCP) Flags() logs.LogWriterFlags {
	return 0
}

func (w *LogWriterGCP) Write(level logs.Level, message string) {
	w.Logger.Log(logging.Entry{
		Severity: levelToGCP(level),
		Payload:  message,
	})
}

func (w *LogWriterGCP) Close() {
	w.Logger.Flush()
	w.Client.Close()
}

func levelToGCP(level logs.Level) logging.Severity {
	switch level {
	case logs.LevelDebug:
		return logging.Debug
	case logs.LevelInfo:
		return logging.Info
	case logs.LevelWarn:
		return logging.Warning
	case logs.LevelError:
		return logging.Error
	case logs.LevelCritical:
		return logging.Critical
	}
	panic("Unknown log level")
}
