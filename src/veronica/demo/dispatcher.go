package demo


import (
    "fmt"
    "time"

    c "veronica/common"
)

type Dispatcher struct {
    c.Dispatcher
    WorkerPool chan (chan int)      // use to get aviable workers
}

func NewDispatcher() *Dispatcher {
    threads := c.Config.Demo.Threads
    workers := []c.WorkerInst{}
    return &Dispatcher{
        Dispatcher : c.NewDispatcher(threads, workers),
        WorkerPool : make(chan (chan int), threads),
    }
}

// Regist and listen to system signal.
func (d *Dispatcher) Init(m <-chan c.SIGNAL) {
    go func() {
        for {
            select {
            case signal := <-m:
                if signal == c.SIGSTART {
                    d.run()
                } else if signal == c.SIGSTOP {
                    d.stop()
                } else if signal == c.SIGMONITOR {
                    d.Monitor()
                }
            case <- time.After(c.DefaultSleepDur):
            }
        }
    }()
}

// Start running dispatcher
func (d *Dispatcher) run() {
    for i := 0; i < d.Threads; i++ {
        id     := i + 1
        worker := NewDemo(id, d.WorkerPool)
        worker.Start()                            // start worker
        d.Workers = append(d.Workers, worker)     // add worker instance to slice
    }
    go d.dispatch()
}

// work dispatch demo
func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <-c.DemoQueue:
            go func(job int) {
                // Try choose one worker to work, if no aviable worker, it will be blocked.
                jobChan := <-d.WorkerPool
                jobChan<- job        // add this job to this worker Chan
            }(job)
        case <- time.After(c.DefaultSleepDur):
        }
    }
}

func (d *Dispatcher) stop() {
    for _, worker := range d.Workers {
        worker.Stop()
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
