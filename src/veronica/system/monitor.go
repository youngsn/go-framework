package system

import (
    "fmt"
    "time"

    c "veronica/common"
)

type Monitor struct {
    name   string
    timer  *time.Ticker
    logger *c.Log
    State  c.RState
}

func NewMonitor() *Monitor {
    interval := c.Config.Monitor.Interval
    logger   := c.NewLog("monitor")
    return &Monitor{
        name   : "Monitor",
        timer  : time.NewTicker(time.Second * time.Duration(interval)),
        logger : logger,
        State  : c.Stopped,
    }
}

func (m *Monitor) Start() {
    m.State = c.Running
    go func() {
        for m.State == c.Running {
            select {
            case <-m.timer.C:
                m.logger.Infof("system")
                m.system()
                m.logger.Infof("module")
                m.module()
            case <-time.After(c.DefaultSleepDur):
            }
        }
    }()
    m.logger.WithFields(c.LogFields{
        "module" : m.name,
    }).Infof("%s, start work", m.name)
}

func (m *Monitor) Stop() {
    m.State = c.Stopped
    m.logger.WithFields(c.LogFields{
        "module" : m.name,
    }).Infof("%s, stopped", m.name)
}

// Module status monitor
func (m *Monitor) module() {
    for _, module := range SysManager.Modules {
        mPacks := module.Monitor()
        for _, pack := range mPacks {
            if pack.Fields == nil {
                pack.Fields = c.LogFields{
                    "module" : pack.Name,
                    "state"  : pack.State,
                }
            } else {
                pack.Fields["module"] = pack.Name
                pack.Fields["state"]  = pack.State
            }
            if pack.Level == c.MONITOR_INFO {
                m.logger.WithFields(pack.Fields).Info(pack.Content)
            } else if pack.Level == c.MONITOR_ERROR {
                m.logger.WithFields(pack.Fields).Error(pack.Content)
            } else {
                m.logger.WithFields(pack.Fields).Info(pack.Content)
            }
        }
    }
}

// App status monitor, put chan monitor
func (m *Monitor) system() {
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

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
