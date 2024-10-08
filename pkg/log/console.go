package log

import "fmt"

type ConsoleWriter struct{}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	message := GetWriterOutputFormat(string(p))

	fmt.Printf("[%s] [%s]: %s\n", message.Level, message.ID, message.Message)
	return len(p), nil
}
