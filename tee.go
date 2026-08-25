package logs

// Tee replicates logs to one or more destinations (typically two or more))
type LogWriterTee struct {
	Targets []LogWriter
}

func (t *LogWriterTee) Write(level Level, message string) {
	for _, target := range t.Targets {
		target.Write(level, message)
	}
}

func (t *LogWriterTee) Close() {
	for _, target := range t.Targets {
		target.Close()
	}
}

// Create a new logger that replicates logs to the given destinations
func Tee(targets ...Log) Log {
	t := &LogWriterTee{}
	for _, target := range targets {
		t.Targets = append(t.Targets, target.LogWriter())
	}
	return NewLogFromWriter(t)
}
