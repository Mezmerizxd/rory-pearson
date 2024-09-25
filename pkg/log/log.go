package log

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"github.com/rs/zerolog"
)

// Instances array
var instances = make(map[string]Log)

// Default directory and file mode for log files.
var logsDirectory = "logs/"
var logsFileMode fs.FileMode = 0755

// Config struct holds the configuration options for the logger.
type Config struct {
	ID            string `json:"id"`            // The name of the id being logged.
	ConsoleOutput bool   `json:"consoleOutput"` // Flag to enable/disable console output.
	FileOutput    bool   `json:"fileOutput"`    // Flag to enable/disable file output.
	StoragePath   string `json:"storagePath"`   // The path to store the log files.
}

// Log interface defines the methods available for logging.
type Log interface {
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Debug() *zerolog.Event
	Close() error
}

// log struct implements the Log interface and holds logger configuration and state.
type log struct {
	config        Config
	writer        zerolog.Logger
	file          *os.File
	consoleWriter *ConsoleWriter
}

// New creates a new logger instance based on the provided configuration.
func New(cfg Config) Log {
	// Check if a logger with the same ID already exists.
	if _, ok := instances[cfg.ID]; ok {
		return instances[cfg.ID]
	}

	var writers []io.Writer
	var file *os.File

	// Check if file output is enabled.
	if cfg.FileOutput {
		// Ensure the storage path exists.
		if _, err := os.Stat(cfg.StoragePath); os.IsNotExist(err) {
			// Create the storage path directory if it doesn't exist.
			if err := os.MkdirAll(cfg.StoragePath, logsFileMode); err != nil {
				print("Error creating storage directory")
				panic(err)
			}
		}

		// Ensure the logs directory exists.
		if _, err := os.Stat(cfg.StoragePath + logsDirectory); os.IsNotExist(err) {
			// Create the logs directory if it doesn't exist.
			if err := os.MkdirAll(cfg.StoragePath+logsDirectory, logsFileMode); err != nil {
				print("Error creating logs directory")
				panic(err)
			}
		}

		// Open or create the log file.
		var err error
		file, err = os.OpenFile(cfg.StoragePath+logsDirectory+cfg.ID+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, logsFileMode)
		if err != nil {
			print("Error opening log file")
			panic(err)
		}
		writers = append(writers, file)
	}

	// If console output is enabled, add os.Stdout to the writers.
	var consoleWriter *ConsoleWriter
	if cfg.ConsoleOutput {
		consoleWriter = NewConsoleWriter()
		writers = append(writers, consoleWriter)
	}

	// Create a multi-writer to output to both file and console if both are enabled.
	multiWriter := io.MultiWriter(writers...)

	// Create a zerolog logger using the multi-writer.
	writer := zerolog.New(multiWriter).With().Timestamp().Logger()

	// Append the new logger to the instances map.
	instances[cfg.ID] = &log{
		config:        cfg,
		writer:        writer,
		file:          file,
		consoleWriter: consoleWriter,
	}

	return instances[cfg.ID]
}

func Get(id string) Log {
	if log, ok := instances[id]; ok {
		return log
	}
	return nil
}

// Info logs an informational message.
func (l *log) Info() *zerolog.Event {
	return l.writer.Info().Str("id", l.config.ID)
}

// Warn logs a warning message.
func (l *log) Warn() *zerolog.Event {
	return l.writer.Warn().Str("id", l.config.ID)
}

// Error logs an error message.
func (l *log) Error() *zerolog.Event {
	return l.writer.Error().Str("id", l.config.ID)
}

// Debug logs a debug message.
func (l *log) Debug() *zerolog.Event {
	return l.writer.Debug().Str("id", l.config.ID)
}

// Close finalizes the logging by closing any open resources, such as file handles.
func (l *log) Close() error {
	l.Info().Msgf("Shutting down %s", l.config.ID)

	// Close the file only if file output is enabled and the file is not nil.
	if l.config.FileOutput && l.file != nil {
		l.file.Close()
	}

	return nil
}

// WriterOutputFormat defines the structure of a log output format.
type WriterOutputFormat struct {
	Level   string `json:"level"`   // Log level (info, warn, error, debug).
	ID      string `json:"id"`      // The id associated with the log message.
	Time    string `json:"time"`    // Timestamp of the log event.
	Message string `json:"message"` // The actual log message.
}

// GetWriterOutputFormat parses JSON log data into the WriterOutputFormat struct.
func GetWriterOutputFormat(data string) WriterOutputFormat {
	var output WriterOutputFormat
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		print("Error unmarshalling log data")
		panic(err)
	}
	return output
}
