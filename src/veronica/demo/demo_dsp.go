package demo


import (
    "time"

    c "veronica/common"
)

type DemoDispatcher struct {
    c.Dispatcher
    WorkPool chan (chan int)      // use to get aviable workers
}

func NewDemoDispatcher() *DemoDispatcher {
    threads := c.Config.Demo.Threads
    workers := []c.WorkerInst{}
    return &DemoDispatcher{
        Dispatcher : c.NewDispatcher(threads, workers),
        WorkPool   : make(chan (chan int), threads),
    }
}

// Regist and listen to system signal.
func (p *DemoDispatcher) Init(m <-chan c.SIGNAL) {
    go func() {
        for {
            select {
            case sg := <-m:
                if sg == c.SIGSTART {
                    p.run()
                } else if sg == c.SIGSTOP {
                    p.stop()
                } else if sg == c.SIGMONITOR {
                    p.Monitor()
                }
            case <- time.After(c.DefaultSleepDur):
            }
        }
    }()
}

// Start running dispatcher
func (p *DemoDispatcher) run() {
    for i := 0; i < p.Threads; i++ {
        id   := i + 1
        inst := NewDemo(id, p.WorkPool)
        inst.Start()                            // start worker
        p.Workers = append(p.Workers, inst)     // add worker instance to slice
    }
    go p.dispatch()
}

// work dispatch demo
func (p *DemoDispatcher) dispatch() {
    for {
        select {
        case job := <-c.DemoQueue:
            go func(job int) {
                // Try choose one worker to work, if no aviable worker, it will be blocked.
                jobChan := <-p.WorkPool
                jobChan<- job        // add this job to this worker Chan
            }(job)
        case <- time.After(c.DefaultSleepDur):
        }
    }
}

func (p *DemoDispatcher) stop() {
    for _, inst := range p.Workers {
        inst.Stop()
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
