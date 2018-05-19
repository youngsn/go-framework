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
    debug int
    addr  string
}

func NewPprof() *PprofMonitor {
    debug := c.Config.Debug.Debug
    addr  := c.Config.Debug.Addr
    return &PprofMonitor{
        debug : debug,
        addr  : addr,
    }
}

// Pprof web monitor
func (p *PprofMonitor) WebMonitor() {
    if p.debug == 0 {           // debug mode will not start pprof
        c.Logger.Infof("pprof debug, off")
        return
    }

    go func() {
        c.Logger.WithFields(c.LogFields{
            "addr" : p.addr,
        }).Debugf("pprof debug, start work")

        if err := http.ListenAndServe(p.addr, nil); err != nil {
            c.Logger.Fatalf(err.Error())
        }
    }()
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
