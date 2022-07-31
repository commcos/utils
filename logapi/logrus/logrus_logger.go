/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

// Package logrus provides an adapter to the
// go-kit log.Logger interface.
package logrus

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xsbull/utils/fileinfo"
	logapi "github.com/xsbull/utils/logapi"
	utilpath "github.com/xsbull/utils/path"
)

//Logger implement
type Logger struct {
	logrus.Ext1FieldLogger
	callerDep        int
	rawLoggerHandler *logrus.Logger

	LocalFile          string
	RotateInterval     int
	MaxFileSize        float64
	outputType         logapi.LoggerOutputType
	initOutputTypeOnce sync.Once
}

const (
	logName = "log"
)

// NewLogger returns a Go kit log.Logger that sends log events to a logrus.Logger.
func NewLogger(callerDep int) logapi.Interface {
	logrusLogger := logrus.New()
	logrusLogger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
		FullTimestamp:   true,
	}
	logrusLogger.SetReportCaller(false)
	logrusLogger.SetLevel(logrus.InfoLevel)

	l := &Logger{
		Ext1FieldLogger:  logrusLogger,
		callerDep:        callerDep,
		rawLoggerHandler: logrusLogger,
	}

	return l
}

// initAriLog init ari log
func (l *Logger) initLogrus() error {
	err := utilpath.MkdirPath(l.LocalFile)
	if err != nil {
		fmt.Errorf("MkdirPath path:%s err:%s", l.LocalFile, err)
		return fmt.Errorf("MkdirPath path:%s err:%s", l.LocalFile, err)
	}

	file, err := os.OpenFile(fmt.Sprintf("%s%s.%s", l.LocalFile, logName, time.Now().Format("2006010215")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		l.rawLoggerHandler.Out = file
	} else {
		fmt.Errorf("Failed to log to file, using default stderr")
		return fmt.Errorf("Failed to log to file, using default stderr")
	}
	return nil
}

// recordLogRotate
func (l Logger) logRotate(logPath string) {
	// 避免操作根目录
	if len(logPath) == 0 {
		return
	}
	fileInfos, err := ioutil.ReadDir(logPath)
	if err != nil {
		fmt.Errorf("scan path %s error", logPath)
	} else {
		var recordLogInfos []os.FileInfo
		for _, fileInfo := range fileInfos {
			if strings.HasPrefix(fileInfo.Name(), logName) {
				recordLogInfos = append(recordLogInfos, fileInfo)
			}
		}

		if len(recordLogInfos) > 0 {
			// 按照指定的日志上限 删除recordLog日志
			fileInfos = fileinfo.SortByModifyTime(fileInfos)
			var currFileSizeTotal float64
			for _, fileInfo := range fileInfos {
				currFileSizeTotal = currFileSizeTotal + float64(fileInfo.Size())/(1024.0*1024.0)
				if currFileSizeTotal > l.MaxFileSize {
					fullFileName := fmt.Sprintf("%s/%s", logPath, fileInfo.Name())
					err = os.Remove(fullFileName)
					if err != nil {
						fmt.Errorf("remove file %s err:%v", fullFileName, err)
					} else {
						fmt.Printf("remove file %s", fullFileName)
					}
				}
			}
		}
	}
}

func (l Logger) caller() (string, string, int, string) {

	var callerStr strings.Builder
	callerStr.Grow(128)

	if pc, file, line, ok := runtime.Caller(l.callerDep); ok {

		funcName := runtime.FuncForPC(pc).Name()
		funcIndex := strings.LastIndex(funcName, "/")

		funcName = funcName[funcIndex+1:]

		fileIdx := strings.LastIndex(file, "/")
		file = file[fileIdx+1:]

		callerStr.WriteString(file)
		callerStr.WriteString(":")
		callerStr.WriteString(strconv.Itoa(line))
		callerStr.WriteString(":")
		callerStr.WriteString(funcName)

		return file, funcName, line, callerStr.String()
	}

	return "caller", "logger", 0, ""
}

//Log wrap level and field
func (l Logger) Logf(level logapi.Level, format string, args ...interface{}) {

	var entry *logrus.Entry

	_, _, _, callerStr := l.caller()
	entry = l.Ext1FieldLogger.WithField("caller", callerStr)

	switch level {
	case logapi.PanicLevel:
		entry.Panicf(format, args...)
	case logapi.FatalLevel:
		entry.Fatalf(format, args...)
	case logapi.ErrorLevel:
		entry.Errorf(format, args...)
	case logapi.WarnLevel:
		entry.Warningf(format, args...)
	case logapi.InfoLevel:
		entry.Infof(format, args...)
	case logapi.DebugLevel:
		entry.Debugf(format, args...)
	case logapi.TraceLevel:
		entry.Tracef(format, args...)
	default:
		entry.Printf(format, args...)
	}
}

//LogF wrap level and field
func (l Logger) LogFf(level logapi.Level, fields logapi.Fields, format string, args ...interface{}) {

	logrusFields := logrus.Fields(fields)

	_, _, _, callerStr := l.caller()
	entry := l.Ext1FieldLogger.WithField("caller", callerStr)

	entry = entry.WithFields(logrusFields)

	switch level {
	case logapi.PanicLevel:
		entry.Panicf(format, args...)
	case logapi.FatalLevel:
		entry.Fatalf(format, args...)
	case logapi.ErrorLevel:
		entry.Errorf(format, args...)
	case logapi.WarnLevel:
		entry.Warningf(format, args...)
	case logapi.InfoLevel:
		entry.Infof(format, args...)
	case logapi.DebugLevel:
		entry.Debugf(format, args...)
	case logapi.TraceLevel:
		entry.Tracef(format, args...)
	default:
		entry.Printf(format, args...)
	}

}

//LogF wrap level and field
func (l Logger) Log(level logapi.Level, args ...interface{}) {

	_, _, _, callerStr := l.caller()
	entry := l.Ext1FieldLogger.WithField("caller", callerStr)

	switch level {
	case logapi.PanicLevel:
		entry.Panicln(args...)
	case logapi.FatalLevel:
		entry.Fatalln(args...)
	case logapi.ErrorLevel:
		entry.Errorln(args...)
	case logapi.WarnLevel:
		entry.Warningln(args...)
	case logapi.InfoLevel:
		entry.Infoln(args...)
	case logapi.DebugLevel:
		entry.Debugln(args...)
	case logapi.TraceLevel:
		entry.Traceln(args...)
	default:
		entry.Println(args...)
	}

}

func (l Logger) Level() logapi.Level {
	lvl := l.rawLoggerHandler.GetLevel()
	outLevel := logapi.InfoLevel

	switch lvl {
	case logrus.DebugLevel:
		outLevel = logapi.DebugLevel
	case logrus.ErrorLevel:
		outLevel = logapi.ErrorLevel
	case logrus.FatalLevel:
		outLevel = logapi.FatalLevel
	case logrus.InfoLevel:
		outLevel = logapi.InfoLevel
	case logrus.TraceLevel:
		outLevel = logapi.TraceLevel
	case logrus.PanicLevel:
		outLevel = logapi.PanicLevel
	case logrus.WarnLevel:
		outLevel = logapi.WarnLevel
	}

	return outLevel
}

func (l Logger) SetLevel(lvl logapi.Level) {

	outLevel := logrus.WarnLevel
	switch lvl {
	case logapi.DebugLevel:
		outLevel = logrus.DebugLevel
	case logapi.ErrorLevel:
		outLevel = logrus.ErrorLevel
	case logapi.FatalLevel:
		outLevel = logrus.FatalLevel
	case logapi.InfoLevel:
		outLevel = logrus.InfoLevel
	case logapi.TraceLevel:
		outLevel = logrus.TraceLevel
	case logapi.PanicLevel:
		outLevel = logrus.PanicLevel
	case logapi.WarnLevel:
		outLevel = logrus.WarnLevel
	default:
		l.Log(logapi.ErrorLevel, "not support this level=%v", lvl)
	}
	l.rawLoggerHandler.SetLevel(outLevel)
}

func (l *Logger) SetLoggerOutputFileConfig(conf logapi.LoggerOutputFileConf) {
	l.LocalFile = conf.LocalFile
	l.RotateInterval = conf.RotateInterval
	l.MaxFileSize = conf.MaxFileSize
}

func (l *Logger) SetLoggerOutputType(outputType logapi.LoggerOutputType) {
	l.outputType = outputType

	if l.outputType == logapi.LoggerOutputTypeFile {
		l.initOutputTypeOnce.Do(func() {
			err := l.initLogrus()
			if err != nil {
				fmt.Errorf("SetLoggerOutputType initLogrus err:%v", err)
				return
			}
			// l.runner.Start()
		})
	}
}
