package common

import (
    "fmt"
    "runtime"
    "strings"

    "github.com/sirupsen/logrus"
)

type contextHook struct {
    levels []logrus.Level
}

func NewContextHook(levels ...logrus.Level) logrus.Hook {
    hook := contextHook{
        levels : levels,
    }
    return &hook
}

// Levels implement levels
func (hook contextHook) Levels() []logrus.Level {
    return logrus.AllLevels
}

// implement logrus fire
// we can add public vals in here
func (hook contextHook) Fire(entry *logrus.Entry) error {
    entry.Data["source"] = findCaller()
    entry.Data["app"]    = APP_NAME
    return nil
}

// find source file & func caller
func findCaller() string {
    skip := 5
    file := ""
    line := 0
    for i := 0; i < 10; i++ {
        file, line = getCaller(skip + i)
        if strings.Contains(file, "sirupsen") == false {
            break
        }
    }
    return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
    _, file, line, ok := runtime.Caller(skip)
    if !ok {
        return "", 0
    }
    return file, line
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
