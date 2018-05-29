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

type worker struct {
    id       int
    name     string
    workPool chan (chan int)
    JobChan  chan int
    State    c.RState
}

func newWorker(id int, workPool chan (chan int)) *worker {
    moduleName := "Demo"
    return &worker{
        id       : id,
        name     : moduleName,
        workPool : workPool,               // pool chan to save worker job chan
        State    : c.Stopped,
        JobChan  : make(chan int),         // receive task from outside world
    }
}

func (w *worker) Start() {
    w.State = c.Running
    c.Logger.WithFields(c.LogFields{
        "module" : w.name,
        "workId" : w.id,
    }).Infof("worker ready")
    go func() {
        // regist worker job chan, means the routine is ready to serve
        w.workPool<- w.JobChan
        for w.State == c.Running {
            select {
            case dq := <-w.JobChan:
                c.Logger.WithFields(c.LogFields{
                    "module" : w.name,
                    "workId" : w.id,
                }).Infof("receive rand num: %d", dq)
                w.workPool<- w.JobChan            // after working, put job chan to worker pool
            case <-time.After(c.DefaultSleepDur):
            }
        }
    } ()
}

func (w *worker) Stop() {
    for w.State == c.Running {
        if len(w.JobChan) > 0 {
            time.Sleep(1 * time.Second)
            continue
        } else {
            w.State = c.Stopped
            break
        }
    }
    c.Logger.WithFields(c.LogFields{
        "module" : w.name,
        "workId" : w.id,
    }).Infof("stopped")
}

func (w *worker) Monitor() *c.MonitorPack {
    level := c.MONITOR_INFO
    if w.State == c.Running {
        level  = c.MONITOR_INFO
    } else {
        level  = c.MONITOR_ERROR
    }
    content   := fmt.Sprint(w.State)
    return &c.MonitorPack{
        Name    : w.name,
        State   : w.State,
        Level   : level,
        Content : content,
        Fields  : c.LogFields{
            "workId" : w.id,
        },
    }
}

func (w *worker) Reload() {
}

func (w *worker) RunStatus() c.RState {
    return w.State
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
