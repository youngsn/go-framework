package sys


// System Signal capture & handle.
// Also Signal will run in main thread util program stopped.
// @AUTHOR tangyang

import (
    // "fmt"
    "os"
    "time"
    "os/signal"
    "syscall"

    . "veronica/common"

    Log "github.com/cihub/seelog"
)


type Signal struct {
    signalChan       chan os.Signal
}


func NewSignal() *Signal {
    signalChan       := make(chan os.Signal)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)  // 监听interrupt & kill

    return &Signal{
        signalChan   : signalChan,
    }
}


func (this *Signal) Start() {
    Log.Infof("System start success")
    for {
        signal    := <-this.signalChan
        if signal == syscall.SIGINT || signal == syscall.SIGTERM {  // stop signal
            Log.Infof("Received signal %v, stop programs", signal)
            SysManager.StopModules()                                // stop模块
            Log.Infof("System stop success, byebye~")
            return
        }

        time.Sleep(DefaultSleepDur)
    }
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
