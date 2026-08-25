package logs

// Adds a prefix to log messages
type LogWriterPrefix struct {
	W      LogWriter
	Prefix string
}

func (w *LogWriterPrefix) Write(level Level, message string) {
	w.W.Write(level, w.Prefix+message)
}

func (w *LogWriterPrefix) Close() {
}

// Return a new logger that adds 'prefix ' (i.e. 'prefix' with a space after it) onto every log message
func NewPrefixLogger(log Log, prefix string) Log {
	return NewPrefixLoggerNoSpace(log, prefix+" ")
}

// Return a new logger that adds 'prefix' onto every log message
func NewPrefixLoggerNoSpace(log Log, prefix string) Log {
	return NewLogFromWriter(&LogWriterPrefix{
		W:      log.LogWriter(),
		Prefix: prefix,
	})
}
