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
    signalChan chan os.Signal
}

func NewSignal() *Signal {
    signalChan := make(chan os.Signal)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)      // listen interrupt & kill
    return &Signal{
        signalChan : signalChan,
    }
}

func (s *Signal) Start() {
    for {
        signal   := <-s.signalChan
        if signal == syscall.SIGINT || signal == syscall.SIGTERM {  // stop signal
            c.Logger.WithFields(c.LogFields{
                "SIGNAL" : signal,
            }).Info("receive stop signal")
            SysManager.StopModules()                                // stop module
            c.Logger.Info("stopped, byebye~")
            return
        }
        time.Sleep(c.DefaultSleepDur)
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
