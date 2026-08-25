package logs

import (
	"fmt"
	"os"
	"sync"
)

// LogWriterFile writes logs to a file and keeps one old file when it gets too large.
type LogWriterFile struct {
	filename string
	maxSize  int64
	file     *os.File
	size     int64
	mu       sync.Mutex
}

// NewLogWriterFile creates a log writer for filename. maxSize is in bytes.
func NewLogWriterFile(filename string, maxSize int64) (*LogWriterFile, error) {
	w := &LogWriterFile{
		filename: filename,
		maxSize:  maxSize,
	}
	if err := w.open(); err != nil {
		return nil, err
	}
	if w.size > w.maxSize {
		if err := w.rollover(); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (w *LogWriterFile) Write(level Level, message string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		if err := w.open(); err != nil {
			return
		}
	}

	n, _ := w.file.Write([]byte(formatMessage(level, LogWriterFlagNeedNewline|LogWriterFlagWantDate|LogWriterFlagWantLevel, message)))
	w.size += int64(n)
	if w.size > w.maxSize {
		_ = w.rollover()
	}
}

func (w *LogWriterFile) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}
}

func (w *LogWriterFile) open() error {
	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("open log file %q: %w", w.filename, err)
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("stat log file %q: %w", w.filename, err)
	}
	w.file = file
	w.size = info.Size()
	return nil
}

func (w *LogWriterFile) rollover() error {
	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		if err != nil {
			return fmt.Errorf("close log file %q: %w", w.filename, err)
		}
	}

	oldFilename := w.filename + ".old"
	if err := os.Remove(oldFilename); err != nil && !os.IsNotExist(err) {
		_ = w.open()
		return fmt.Errorf("remove old log file %q: %w", oldFilename, err)
	}
	if err := os.Rename(w.filename, oldFilename); err != nil {
		_ = w.open()
		return fmt.Errorf("rename log file %q to %q: %w", w.filename, oldFilename, err)
	}
	if err := w.open(); err != nil {
		return err
	}
	return nil
}
