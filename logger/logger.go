/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package logger

import (
	logapi "github.com/xsbull/utils/logapi"
	loggerconfig "github.com/xsbull/utils/logapi/loggerbackend"
	loggerfactory "github.com/xsbull/utils/logapi/loggerbackend/factory"
)

type loggerImpl struct {
	level Level
	logapi.Interface
	outputType LoggerOutputType
}

func (lt *loggerImpl) initBackend() {
	conf := &loggerconfig.Config{
		CallerDepth: 3,
	}
	lt.outputType = LoggerOutputTypeStdOut
	lt.Interface = loggerfactory.CreateLogHandle(conf)
}

var defaultLoggerHandler loggerImpl

func init() {
	defaultLoggerHandler.initBackend()
}

func LogLevel() Level {
	return defaultLoggerHandler.level
}

//InitLogger init logger
func SetLogLevel(level Level) {

	defaultLoggerHandler.level = level

	defaultLoggerHandler.SetLevel(level)
}

//Log wrap level and field
func Log(level Level, args ...interface{}) {
	defaultLoggerHandler.Log(level, args...)
}

//LogF wrap level and field
func Logf(level Level, template string, args ...interface{}) {

	defaultLoggerHandler.Logf(level, template, args...)
}

func SetLoggerOutputType(outputType LoggerOutputType) {
	defaultLoggerHandler.outputType = outputType
	defaultLoggerHandler.SetLoggerOutputType(defaultLoggerHandler.outputType)
}

func SetLoggerOutputConfig(conf LoggerOutputFileConf) {
	defaultLoggerHandler.SetLoggerOutputFileConfig(conf)
}
