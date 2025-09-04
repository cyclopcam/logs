package logs

// PrefixLogWriter writes to the underlying log, but all messages are prefixed with a string of your choice
type PrefixLogWriter struct {
	Base   LogWriter
	Prefix string
}

// Create a new PrefixLogWriter
func NewPrefixLogWriter(log Log, prefix string) Log {
	return NewPrefixLogWriterNoSpace(log, prefix+" ")
}

// Create a new PrefixLogWriter, but don't add a space onto 'prefix'
func NewPrefixLogWriterNoSpace(log Log, prefix string) Log {
	writer := &PrefixLogWriter{
		Base:   log.LogWriter(),
		Prefix: prefix,
	}
	return &Logger{
		Output: writer,
	}
}

func (p *PrefixLogWriter) Flags() LogWriterFlags {
	return p.Base.Flags()
}
func (p *PrefixLogWriter) Write(level Level, message string) {
	p.Base.Write(level, p.Prefix+message)
}
func (p *PrefixLogWriter) Close() {
	p.Base.Close()
}
