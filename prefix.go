package logs

// Return a new logger that adds 'prefix ' (i.e. 'prefix' with a space after it) onto every log message
func NewPrefixLogger(log Log, prefix string) Log {
	return NewPrefixLoggerNoSpace(log, prefix+" ")
}

// Return a new logger that adds 'prefix' onto every log message
func NewPrefixLoggerNoSpace(log Log, prefix string) Log {
	logger, ok := log.(*Logger)
	if !ok {
		panic("Underlying log is not a Logger")
	}
	newLogger := &Logger{
		Output: logger.Output,
		Prefix: logger.Prefix + prefix,
	}
	return newLogger
}
