# Logs

This is the logging package used by the Cyclops camera system.

These are plain old textual logs, with a logging level and a format string.

You can write your own custom log output object by implementing the LogWriter interface,
and replacing the Output object on your Log. For example, if you want to store log
messages in a buffer, then you can do this:

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
