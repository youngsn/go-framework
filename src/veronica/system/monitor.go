package system


import (
    "time"

    . "veronica/common"

    "github.com/cihub/seelog"
)


type Monitor struct {
    State           RState
    saveQueue       chan *MonitorPack

    timer           *time.Ticker
    logger          seelog.LoggerInterface
}


func NewMonitor() *Monitor {
    interval        := Config.Global.MonitorInterval
    logger          := GetLogger("monitor")

    return &Monitor{
        State       : Stopped,
        timer       : time.NewTicker(time.Second * time.Duration(interval)),
        logger      : logger,
    }
}


func (this *Monitor) Start() {
    this.State      = Running
    go this.run()
}


func (this *Monitor) run() {
    this.logger.Infof("Monitor thread, started")
    for this.State == Running {
        select {
        case <-this.timer.C:
            this.logger.Infof("Program status @%s", time.Now().Format("2006-01-02 15:04:05"))
            this.logger.Infof("System:")
            this.systemMonitor()

            this.logger.Infof("Modules:")
            this.modulesMonitor()
        default:
            time.Sleep(DefaultSleepDur)
        }
    }
}


func (this *Monitor) Stop() {
    this.State       = Stopped
    this.logger.Infof("Monitor thread, stopped")
}


func (this *Monitor) modulesMonitor() {
    for moduleName, module := range SysManager.Modules {
        this.logger.Infof("%s Status:", moduleName)

        monitorPacks       := module.Monitor()
        for _, monitorPack := range monitorPacks {
            if monitorPack.StdLevel == MONITOR_INFO {
                this.logger.Infof(monitorPack.Content)
            } else if monitorPack.StdLevel == MONITOR_WARN {
                this.logger.Warnf(monitorPack.Content)
            } else if monitorPack.StdLevel == MONITOR_ERROR {
                this.logger.Errorf(monitorPack.Content)
            } else if monitorPack.StdLevel == MONITOR_FATAL {
                this.logger.Criticalf(monitorPack.Content)
            } else {
                this.logger.Infof(monitorPack.Content)
            }
        }
    }
}


// 内部常量监控，etc：内部队列情况
func (this *Monitor) systemMonitor() {
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
