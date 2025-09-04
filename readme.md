# Logs

This is the logging package used by the Cyclops camera system.

These are plain old textual logs, with a logging level and a format string.

You can write your own custom log output object by implementing the LogWriter interface,
and replacing the Output object on your Log. For example, if you want to store log
messages in a buffer, then you can do this:

### Base:

To use:
`go get github.com/cyclopcam/logs/v3@v3.0.2`

```go
import "github.com/cyclopcam/logs/v3"

func foo() {
	logs.NewLog()
}
```

### For GCP:

To use inside a GCP project, also do this:
`go get github.com/cyclopcam/logs/gcp/v3@v3.0.2`

```go
import logsgcp "github.com/cyclopcam/logs/gcp/v3"

func foo() {
	logsgcp.NewLog()
}
```

```go

// LogStore is a log writer that stores log messages in a slice of strings before sending them out
type LogStore struct {
	Stored []string       // Stored logs
	Output logs.LogWriter // Original writer
}

func (s *LogStore) Flags() logs.LogWriterFlags {
	return s.Output.Flags()
}

func (s *LogStore) Write(level logs.Level, message string) {
	s.Stored = append(s.Stored, message)
	s.Output.Write(level, message)
}

func (s *LogStore) Close() {
	// TODO: Do something special with the stored messages, like send them up to a telemetry server
	s.Output.Close()
}

```

## Multi-module repo

This is a multi-module git repo, meaning there are different go.mod files inside here.

To tag a new release of the base package:

`git tag v2.0.1`

To tag a new release of the child (GCP) package:

`git tag gcp/v2.0.5`

The `gcp` prefix tells Go that this tag is only for the `gcp` module within the repo.

To work inside this repo, you might want to do this:

`go work init ./ ./gcp`
