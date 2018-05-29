package system


// System Signal capture & handle.
// Also Signal will run in main thread util program stopped.
// @AUTHOR tangyang
import (
    "os"
    "time"
    "syscall"
    "os/signal"

    c "veronica/common"
)

type Signal struct {
    ch chan os.Signal
}

func NewSignal() *Signal {
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)      // listen interrupt & kill
    return &Signal{
        ch : ch,
    }
}

func (p *Signal) Start() {
    for {
        sg := <-p.ch
        if sg == syscall.SIGINT || sg == syscall.SIGTERM {  // stop signal
            c.Logger.WithFields(c.LogFields{
                "SIGNAL" : sg,
            }).Infof("received signal, %s", sg)
            SysManager.Stop()                               // stop module
            c.Logger.Infof("success, bye~")
            return
        }
        time.Sleep(c.DefaultSleepDur)
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
