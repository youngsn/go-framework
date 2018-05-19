package common

// All util functions are here.
// @author tangyang
import (
    "os"
    "fmt"
    "syscall"
    "strconv"
)

func FilePathExist(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        if err := os.Mkdir(path, 0775); err != nil {
            return err
        }
    }
    return nil
}

// Write pid into pidfile in var dir
func WritePid() error {
    path    := RunPath + "/" + APP_NAME + ".pid"
    pid     := os.Getpid()
    fd, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    defer fd.Close()

    if err != nil {
        fmt.Sprintf("file %s, %s", path, err.Error())
    }
    fd.Write([]byte(strconv.Itoa(pid)))
    return nil
}

func UnlinkPid() {
    pidFile := RunPath + "/" + APP_NAME + ".pid"
    syscall.Unlink(pidFile)
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
