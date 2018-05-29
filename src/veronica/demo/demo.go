package demo


import (
    "time"

    c "veronica/common"
)

type Demo struct {
    c.Worker
    workPool chan (chan int)
    JobChan  chan int
}

func NewDemo(id int, workPool chan (chan int)) *Demo {
    name := "Demo"
    return &Demo{
        Worker   : c.NewWorker(id, name),
        workPool : workPool,               // pool chan to save worker job chan
        JobChan  : make(chan int),         // receive task from outside world
    }
}

func (p *Demo) Start() {
    p.State = c.Running
    c.Logger.WithFields(c.LogFields{
        "module" : p.Name,
        "workId" : p.Id,
    }).Infof("worker ready")
    go func() {
        // regist worker job chan, means the routine is ready to serve
        p.workPool<- p.JobChan
        for p.State == c.Running {
            select {
            case dq := <-p.JobChan:
                c.Logger.WithFields(c.LogFields{
                    "module" : p.Name,
                    "workId" : p.Id,
                }).Infof("receive rand num: %d", dq)
                p.workPool<- p.JobChan            // after working, put job chan to worker pool
            case <-time.After(c.DefaultSleepDur):
            }
        }
    }()
}

func (p *Demo) Stop() {
    for p.State == c.Running {
        if len(p.JobChan) > 0 {
            time.Sleep(1 * time.Second)
            continue
        } else {
            p.State = c.Stopped
            break
        }
    }
    c.Logger.WithFields(c.LogFields{
        "module" : p.Name,
        "workId" : p.Id,
    }).Infof("stopped")
}

func (p *Demo) Monitor() *c.MonitorPack {
    level := c.MONITOR_INFO
    if p.State == c.Running {
        level  = c.MONITOR_INFO
    } else {
        level  = c.MONITOR_ERROR
    }
    content   := p.State.String()
    return &c.MonitorPack{
        Name    : p.Name,
        State   : p.State,
        Level   : level,
        Content : content,
        Fields  : c.LogFields{
            "workId" : p.Id,
        },
    }
}
func (p *Demo) Reload() {
}

func (p *Demo) RunStatus() c.RState {
    return p.State
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
