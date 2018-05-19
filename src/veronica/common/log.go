package common

import (
    "fmt"
    "path"
    "time"

    rotatelog "github.com/lestrrat-go/file-rotatelogs"    // log rotate
    lfs "github.com/rifflock/lfshook"                     // log to file
    "github.com/sirupsen/logrus"                          // log engine
)

type LogFields map[string]interface{}

type Log struct {
    log *logrus.Logger         // log engine
}

func NewLog(name string) *Log {
    logConf   := Config.Log

    // init log formatter
    formatter := new(logrus.TextFormatter)
    formatter.TimestampFormat = "2006-01-02 15:04:05"
    formatter.DisableSorting  = false

    // init file log hook, file log using lfs hook
    path, wfPath := getLogFilepath(name, logConf.Path)  // log paths
    hook      := buildLogHook(path, wfPath, logConf.Rotate, formatter)

    log       := logrus.New()
    log.Level  = getLogLevel(logConf.Level)             // log level
    log.AddHook(hook)
    log.AddHook(NewContextHook())                       // public context hook
    return &Log{
        log : log,
    }
}

func buildLogHook(path, wfPath string, rotate int, formatter logrus.Formatter) *lfs.LfsHook {
    if 1 == rotate {
        pathWriter, err := rotatelog.New(
            path + ".%Y%m%d%H",
            rotatelog.WithLinkName(path),
            rotatelog.WithRotationTime(time.Hour),
        )
        if err != nil {
            panic(fmt.Sprintf("log error, %s", err.Error()))
        }

        wfPathWriter, err := rotatelog.New(
            wfPath + ".%Y%m%d%H",
            rotatelog.WithLinkName(wfPath),
            rotatelog.WithRotationTime(time.Hour),
        )
        if err != nil {
            panic(fmt.Sprintf("log error, %s", err.Error()))
        }

        lfsMap := lfs.WriterMap{
            logrus.DebugLevel : pathWriter,
            logrus.InfoLevel  : pathWriter,
            logrus.ErrorLevel : wfPathWriter,
            logrus.FatalLevel : wfPathWriter,
            logrus.PanicLevel : wfPathWriter,
        }
        return lfs.NewHook(
            lfsMap,
            formatter,
        )
    } else {
        lfsMap := lfs.PathMap{
            logrus.DebugLevel : path,
            logrus.InfoLevel  : path,
            logrus.ErrorLevel : wfPath,
            logrus.FatalLevel : wfPath,
            logrus.PanicLevel : wfPath,
        }
        return lfs.NewHook(
            lfsMap,
            formatter,
        )
    }
}

// get log filepath
func getLogFilepath(name, filepath string) (string, string) {
    if path.IsAbs(filepath) == false {                      // not abspath
        filepath = fmt.Sprintf("%s/%s", SysPath, filepath)
    }
    if err := FilePathExist(filepath); err != nil {
        panic(fmt.Sprintf("log filepath error, path: %s, %s", filepath, err.Error()))
    }

    var filename string
    if filename  = fmt.Sprintf("worker-%s.log", name); name == "app" {
        filename = "worker.log"
    }

    path   := fmt.Sprintf("%s/%s", filepath, filename)      // debug|info log path
    wfpath := fmt.Sprintf("%s/%s.wf", filepath, filename)   // error|fatal|panic log path
    return path, wfpath
}

// logrus level map
func getLogLevel(level int) logrus.Level {
    logLevel := logrus.InfoLevel
    if 1 == level {
        logLevel = logrus.DebugLevel
    } else if 2 == level {
        logLevel = logrus.InfoLevel
    } else if 4 == level {
        logLevel = logrus.ErrorLevel
    } else if 8 == level {
        logLevel = logrus.FatalLevel
    } else if 16 == level {
        logLevel = logrus.PanicLevel
    }
    return logLevel
}

func (l *Log) WithFields(fields LogFields) *logrus.Entry {
    return l.log.WithFields(logrus.Fields(fields))
}

func (l *Log) Debug(args ...interface{}) {
    l.log.Debug(args)
}

func (l *Log) Debugf(format string, args ...interface{}) {
    l.log.Debugf(format, args...)
}

func (l *Log) Info(args ...interface{}) {
    l.log.Info(args)
}

func (l *Log) Infof(format string, args ...interface{}) {
    l.log.Infof(format, args...)
}

func (l *Log) Error(args ...interface{}) {
    l.log.Error(args)
}

func (l *Log) Errorf(format string, args ...interface{}) {
    l.log.Errorf(format, args...)
}

func (l *Log) Fatal(args ...interface{}) {
    l.log.Fatal(args)
}

func (l *Log) Fatalf(format string, args ...interface{}) {
    l.log.Fatalf(format, args...)
}

func (l *Log) Panic(args ...interface{}) {
    l.log.Panic(args)
}

func (l *Log) Panicf(format string, args ...interface{}) {
    l.log.Panicf(format, args...)
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
