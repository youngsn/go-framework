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
    State  c.RState
}

func NewMonitor() *Monitor {
    interval := c.Config.Monitor.Interval
    logger   := c.NewLog("monitor")
    return &Monitor{
        name   : "Monitor",
        tk     : time.NewTicker(time.Second * time.Duration(interval)),
        logger : logger,
        State  : c.Stopped,
    }
}

func (m *Monitor) run() {
    m.State = c.Running
    go func() {
        for m.State == c.Running {
            select {
            case <-m.tk.C:
                m.logger.Infof("system")
                m.systemMonitor()
                SysManager.SendBoardcast(c.SIGMONITOR)      // send monitor sig
            case packs := <-c.MonitorQueue:
                m.logger.Infof("module")
                m.moduleMonitor(packs)
            case <-time.After(c.DefaultSleepDur):
            }
        }
    }()
}

func (m *Monitor) stop() {
    m.State = c.Stopped
}

// Output module monitor info
func (m *Monitor) moduleMonitor(mPacks []*c.MonitorPack) {
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
            m.logger.WithFields(pack.Fields).Errorf(pack.Content)
        } else {
            m.logger.WithFields(pack.Fields).Info(pack.Content)
        }
    }
}

// App status monitor, put chan monitor
func (m *Monitor) systemMonitor() {
    chs := map[string]struct {Len int; Cap int}{
         "DemoQueue" : {Len : len(c.DemoQueue),  Cap : cap(c.DemoQueue)},
    }
    for name, ch := range chs {
        m.logger.WithFields(c.LogFields{
            "len" : ch.Len,
            "rat" : fmt.Sprintf("%d/%d", ch.Len, ch.Cap),
        }).Infof(name)
    }
}

func (m *Monitor) Init(ch <-chan c.SIGNAL) {
    go func() {
        for {
            select {
            case sg := <-ch:
                if sg == c.SIGSTART {
                    m.run()
                } else if sg == c.SIGSTOP {
                    m.stop()
                }
            case <- time.After(c.DefaultSleepDur):
            }
        }
    }()
}

func (m *Monitor) Status() c.RState {
    return m.State
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
