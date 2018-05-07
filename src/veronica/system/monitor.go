package system

import (
    "time"

    c "veronica/common"
)

type Monitor struct {
    State  c.RState
    timer  *time.Ticker
    logger *c.Log
}

func NewMonitor() *Monitor {
    interval := c.Config.Monitor.Interval
    logger   := c.NewLog("monitor")
    return &Monitor{
        State  : c.Stopped,
        timer  : time.NewTicker(time.Second * time.Duration(interval)),
        logger : logger,
    }
}

func (m *Monitor) Start() {
    m.State = c.Running
    go m.run()
    m.logger.WithFields(c.LogFields{
        "moduleName" : "monitor",
    }).Info("monitor start")
}

func (m *Monitor) run() {
    for m.State == c.Running {
        select {
        case <-m.timer.C:
            m.logger.Info("system")
            m.systemMonitor()

            m.logger.Info("module")
            m.modulesMonitor()
        default:
            time.Sleep(c.DefaultSleepDur)
        }
    }
}

func (m *Monitor) Stop() {
    m.State = c.Stopped
    m.logger.WithFields(c.LogFields{
        "moduleName" : "monitor",
    }).Infof("monitor stopped")
}

func (m *Monitor) modulesMonitor() {
    for _, module := range SysManager.Modules {
        monitorPacks := module.Monitor()
        for _, mpack := range monitorPacks {
            if mpack.StdLevel == c.MONITOR_INFO {
                m.logger.WithFields(mpack.Fields).Info(mpack.Content)
            } else if mpack.StdLevel == c.MONITOR_ERROR {
                m.logger.WithFields(mpack.Fields).Error(mpack.Content)
            } else if mpack.StdLevel == c.MONITOR_FATAL {
                m.logger.WithFields(mpack.Fields).Fatal(mpack.Content)
            } else if mpack.StdLevel == c.MONITOR_PANIC {
                m.logger.WithFields(mpack.Fields).Panic(mpack.Content)
            } else {
                m.logger.WithFields(mpack.Fields).Info(mpack.Content)
            }
        }
    }
}

// system monitors can set here
func (this *Monitor) systemMonitor() {
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
