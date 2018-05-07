package common

// All util functions are here.
// @author tangyang
import (
    "os"
    "fmt"
    "syscall"
    "strconv"
)

// If filepath exists, will auto create one if not exist.
func FilePathExist(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        if err := os.Mkdir(path, 0775); err != nil {
            return err
        }
    }
    return nil
}

// Write main pid into pidfile.
func WritePid() error {
    pidFile := RunPath + "/" + APP_NAME + ".pid"
    pid     := os.Getpid()
    fd, err := os.OpenFile(pidFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    defer fd.Close()

    if err != nil {
        fmt.Sprintf("Can not open pid file %s, %s", pidFile, err.Error())
    }
    fd.Write([]byte(strconv.Itoa(pid)))
    return nil
}

// Unlink pid file
func UnlinkPid() {
    pidFile := RunPath + "/" + APP_NAME + ".pid"
    syscall.Unlink(pidFile)
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
