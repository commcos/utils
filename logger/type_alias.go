package logger

import (
	logapi "github.com/commcos/utils/logapi"
)

type LoggerOutputFileConf = logapi.LoggerOutputFileConf

type LoggerOutputType = logapi.LoggerOutputType

const (
	LoggerOutputTypeStdOut = logapi.LoggerOutputTypeStdOut
	LoggerOutputTypeFile   = logapi.LoggerOutputTypeFile
)

// Fields log fields
type Fields = logapi.Fields

// Level logger level
type Level = logapi.Level

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = logapi.PanicLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = logapi.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = logapi.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = logapi.WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = logapi.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = logapi.DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = logapi.TraceLevel
)
