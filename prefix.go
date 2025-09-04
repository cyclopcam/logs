package logs

// PrefixLogger writes to the underlying log, but all messages are prefixed with a string of your choice
type PrefixLogger struct {
	Base   LogWriter
	Prefix string
}

// Create a new PrefixLogWriter
func NewPrefixLogger(log Log, prefix string) Log {
	return NewPrefixLoggerNoSpace(log, prefix+" ")
}

// Create a new PrefixLogWriter, but don't add a space onto 'prefix'
func NewPrefixLoggerNoSpace(log Log, prefix string) Log {
	writer := &PrefixLogger{
		Base:   log.LogWriter(),
		Prefix: prefix,
	}
	return &Logger{
		Output: writer,
	}
}

func (p *PrefixLogger) Flags() LogWriterFlags {
	return p.Base.Flags()
}
func (p *PrefixLogger) Write(level Level, message string) {
	p.Base.Write(level, p.Prefix+message)
}
func (p *PrefixLogger) Close() {
	p.Base.Close()
}
