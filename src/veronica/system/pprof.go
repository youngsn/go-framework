package system

// Pprof monitor.
// Very useful when program is debugging
// Be careful if use in production.
// You can config start it or not in config file.
// @AUTHOR tangyang
import (
    "net/http"
    _ "net/http/pprof"
    c "veronica/common"
)

type PprofMonitor struct {
    debug  int
    remote string
}

func NewPprof() *PprofMonitor {
    debug  := c.Config.Debug.Debug
    remote := c.Config.Debug.Remote
    return &PprofMonitor{
        debug  : debug,
        remote : remote,
    }
}

// Pprof web. 
func (p *PprofMonitor) WebMonitor() {
    if p.debug == 0 {           // debug mode will not start pprof
        return
    }

    go func() {
        c.Logger.WithFields(c.LogFields{
            "addr" : p.remote,
        }).Info("start debug mode")

        if err := http.ListenAndServe(p.remote, nil); err != nil {
            c.Logger.Fatalf(err.Error())
        }
    }()
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
