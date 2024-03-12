package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// https://github.com/rs/zerolog?tab=readme-ov-file#pretty-logging
var consoleWriter io.Writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}

type safeConsoleWriter struct {
	lock   sync.Mutex
	writer io.Writer
}

func newSafeConsoleWriter() io.Writer {
	return &safeConsoleWriter{
		writer: consoleWriter,
	}
}

func (s *safeConsoleWriter) Write(p []byte) (n int, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.writer.Write(p)
}
