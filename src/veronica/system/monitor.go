package system

import (
    "fmt"
    "time"

    c "veronica/common"
)

type Monitor struct {
    name   string
    tk     *time.Ticker
    logger *c.Log
    state  c.RState
}

func NewMonitor() *Monitor {
    interval := c.Config.Monitor.Interval
    logger   := c.NewLog("monitor")
    return &Monitor{
        name   : "Monitor",
        tk     : time.NewTicker(time.Second * time.Duration(interval)),
        logger : logger,
        state  : c.Stopped,
    }
}

func (p *Monitor) run() {
    p.state = c.Running
    go func() {
        for p.state == c.Running {
            select {
            case <-p.tk.C:
                p.systemMonitor()
                SysManager.SendBoardcast(c.SIGMONITOR)      // send monitor sig
            case packs := <-c.MonitorQueue:
                p.moduleMonitor(packs)
            case <-time.After(c.DefaultSleepDur):
            }
        }
    }()
}

func (p *Monitor) stop() {
    p.state = c.Stopped
}

// Output module monitor info
func (p *Monitor) moduleMonitor(mPacks []*c.MonitorPack) {
    if mPacks == nil {
        return
    }
    for _, pack := range mPacks {
        name    := pack.Name
        if pack.Fields == nil {
            pack.Fields = c.LogFields{
                "module" : name,
                "state"  : pack.State,
            }
        } else {
            pack.Fields["module"] = name
            pack.Fields["state"]  = pack.State
        }
        if pack.Level == c.MONITOR_ERROR {
            p.logger.WithFields(pack.Fields).Errorf(pack.Content)
        } else {
            p.logger.WithFields(pack.Fields).Info(pack.Content)
        }
    }
}

// App status monitor, put chan monitor
func (p *Monitor) systemMonitor() {
    chs := map[string]struct {Len int; Cap int}{
         "DemoQueue" : {Len : len(c.DemoQueue),  Cap : cap(c.DemoQueue)},
    }
    for name, ch := range chs {
        p.logger.WithFields(c.LogFields{
            "len" : ch.Len,
            "rat" : fmt.Sprintf("%d/%d", ch.Len, ch.Cap),
        }).Infof(name)
    }
}

func (p *Monitor) Init(ch <-chan c.SIGNAL) {
    go func() {
        for {
            select {
            case sg := <-ch:
                if sg == c.SIGSTART {
                    p.run()
                } else if sg == c.SIGSTOP {
                    p.stop()
                }
            case <- time.After(c.DefaultSleepDur):
            }
        }
    }()
}

func (p *Monitor) Status() c.RState {
    return p.state
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
