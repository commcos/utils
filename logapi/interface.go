/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package logapi

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

//Level logger level
type Level uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

//String2Level convert string level into level
func String2Level(level string) Level {
	resultLevel := InfoLevel

	switch level {
	case "trace":
		resultLevel = TraceLevel
	case "debug":
		resultLevel = DebugLevel
	case "info":
		resultLevel = InfoLevel
	case "warn":
		resultLevel = WarnLevel
	case "error":
		resultLevel = ErrorLevel
	case "fatal":
		resultLevel = FatalLevel
	case "panic":
		resultLevel = PanicLevel
	}

	return resultLevel
}

type LoggerOutputType string

const (
	LoggerOutputTypeStdOut LoggerOutputType = "stdout"
	LoggerOutputTypeFile   LoggerOutputType = "file"
)

type LoggerOutputFileConf struct {
	LocalFile      string
	RotateInterval int
	MaxFileSize    float64
}

//Interface extend Logger support. add level with Logger
type Interface interface {
	Logf(level Level, format string, args ...interface{})
	Log(level Level, args ...interface{})
	Level() Level
	SetLevel(Level)
	SetLoggerOutputFileConfig(conf LoggerOutputFileConf)
	SetLoggerOutputType(outputType LoggerOutputType)
}
