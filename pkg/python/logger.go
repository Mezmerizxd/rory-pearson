package python

import (
	"fmt"
	"rory-pearson/pkg/log"
)

type LoggerWriter struct {
	Log  log.Log
	Type string
}

func NewLoggerWriter(c LoggerWriter) *LoggerWriter {
	return &LoggerWriter{
		Log:  c.Log,
		Type: c.Type,
	}
}

func (w *LoggerWriter) Write(p []byte) (n int, err error) {
	switch w.Type {
	case "info":
		w.Log.Info().Msg(string(p))
	case "error":
		w.Log.Error().Msg(string(p))
	default:
		w.Log.Warn().Msg(fmt.Sprintf("[UNKNOWN] %s", p))
	}

	return len(p), nil
}
