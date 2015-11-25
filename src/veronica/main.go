package main


// Program start here.
// @AUTHOR tangyang

import (
   "os"
   "runtime"

    s "veronica/system"
    . "veronica/common"
)


func main() {
    if err := s.Initialize(); err != nil {
        panic(err.Error())
    }

    if err := WritePid(); err != nil {
        panic(err.Error())
    }

    processors    := Config.Global.MaxProcs
    runtime.GOMAXPROCS(processors)

    s.SysPprofMonitor.WebPprofMonitor()     // pprof
    if err := s.SysManager.StartModules(); err != nil {             // modules
        panic(err.Error())
    }

    sysSignal     := s.NewSignal()          // system signal
    sysSignal.Start()

    UnlinkPid()
    os.Exit(0)
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
