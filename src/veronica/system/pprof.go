package sys


// Pprof monitor.
// Very useful when program is debugging
// Be careful if use in production.
// You can config start it or not in config file.
// @AUTHOR tangyang

import(
    "net/http"
    _ "net/http/pprof"

    . "veronica/common"

    Log "github.com/cihub/seelog"
)


type PprofMonitor struct {
    debug         bool
    addr          string
}


func NewPprofMonitor() *PprofMonitor {
    debug        := Config.Global.PprofMode
    addr         := Config.Global.PprofAddr

    return &PprofMonitor{
        debug    : debug,
        addr     : addr,
    }
}


// Pprof web. 
func (this *PprofMonitor) WebPprofMonitor() {
    if this.debug == false {
        Log.Infof("Not Debug mode, PprofMonitor not started")
        return
    }

    go func(){
        Log.Warnf("Start Web PprofMonitor")
        if err := http.ListenAndServe(this.addr, nil); err != nil {
            Log.Errorf(err.Error())
        }
    }()
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
