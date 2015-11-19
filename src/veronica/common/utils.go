package common

// All tool funcs belongs here.
// @author tangyang

import (
    "os"
    "fmt"
    "syscall"
    "strconv"

    "github.com/cihub/seelog"
)


// Get seelog instance according to module name.
func GetLogger(loggerName string) seelog.LoggerInterface {
    if logInstance, ok := LoggerFactory["default"]; ok {
        return logInstance
    } else {
        panic(fmt.Sprintf("logger %s not exist", loggerName))
    }
}


// If filepath exists, will auto create one if not exist.
func FileExist(path string) error {
    if _, err  := os.Stat(path); os.IsNotExist(err) {
        if err := os.Mkdir(path, 0775); err != nil {
            return err
        }
    }

    return nil
}


// Write main pid into pidfile.
var pidFile string    = RunPath + "run.pid"      // pid
func WritePid() error {
    pid              := os.Getpid()
    fd, err          := os.OpenFile(pidFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    defer fd.Close()

    if err != nil {
        fmt.Sprintf("Can not open pid file %s, %s", pidFile, err.Error())
    }
    fd.Write([]byte(strconv.Itoa(pid)))

    return nil
}


func UnlinkPid() {
    syscall.Unlink(pidFile)
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
