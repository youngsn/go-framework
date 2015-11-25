package demo


import (
    "fmt"
    "time"

    . "veronica/common"
)


type DManager struct {
    pThreads       int
    Demos          []*Demo
}


func NewDManager() *DManager {
    pThreads       := Config.Demo.Threads
    return &DManager{
        pThreads   : pThreads,
        Demos      : []*Demo{},
    }
}


func (this *DManager) Status() RState {
    for _, demo := range this.Demos {
        if demo.State == Running {
            return Running
        } else {
            continue
        }
    }

    return Stopped
}


func (this *DManager) Monitor() []*MonitorPack {
    monitorStatus := []*MonitorPack{}
    for _, demo := range this.Demos {
        monitorStatus = append(monitorStatus, demo.ProcStatus())
    }

    return monitorStatus
}


func (this *DManager) Ctrl(m <-chan SIGNAL) {
    go func() {
        for {
            select {
            case signal := <-m:
                if signal == SIGSTART {
                    this.run()
                } else if signal == SIGSTOP {
                    this.stop()
                }
            default:
                time.Sleep(DefaultSleepDur)
            }
        }
    }()
}


func (this *DManager) run() {
    for i := 0; i < this.pThreads; i++ {
        id               := i + 1
        demo             := NewDemo(id)
        demo.Start()

        this.Demos        = append(this.Demos, demo)
    }
}


func (this *DManager) stop() {
    for _, demo := range this.Demos {
        demo.Stop()
    }
}


type Demo struct {
    Id            int
    State         RState
}


func NewDemo(id int) *Demo {
    return &Demo{
        Id          : id,
        State       : Stopped,
    }
}


func (this *Demo) Start() {
    this.State        = Running
    go this.run()
}


func (this *Demo) run() {
    Log.Infof("Thread[%d], start", this.Id)

    for this.State == Running {
        select {
        case dq := <-DemoQueue:
            Log.Infof("Thread[%d], Demo received %v", this.Id, dq)
        default:
            time.Sleep(DefaultSleepDur)
        }
    }
}


func (this *Demo) Stop() {
    for this.State == Running {
        if len(DemoQueue) > 0 {
            time.Sleep(1 * time.Second)
            continue
        } else {
            this.State             = Stopped
            break
        }
    }

    Log.Infof("Thread[%d], stop success", this.Id)
}


func (this *Demo) ProcStatus() *MonitorPack {
    stdLevel         := MONITOR_INFO
    if this.State == Running {
        stdLevel      = MONITOR_INFO
    } else {
        stdLevel      = MONITOR_ERROR
    }

    content          := fmt.Sprintf("Thread[%d], state %s", this.Id, this.State)
    return &MonitorPack{
        StdLevel     : stdLevel,
        Content      : content,
    }
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
