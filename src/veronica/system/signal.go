package system


// System Signal capture & handle.
// Also Signal will run in main thread util program stopped.
// @AUTHOR tangyang
import (
    "os"
    "time"
    "os/signal"
    "syscall"

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

func (s *Signal) Start() {
    for {
        signal   := <-s.ch
        if signal == syscall.SIGINT || signal == syscall.SIGTERM {  // stop signal
            c.Logger.WithFields(c.LogFields{
                "SIGNAL" : signal,
            }).Infof("received signal, %s", signal)
            SysManager.Stop()                                       // stop module
            c.Logger.Infof("success, bye~")
            return
        }
        time.Sleep(c.DefaultSleepDur)
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
