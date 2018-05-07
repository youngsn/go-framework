package demo


import (
    "fmt"
    "time"

    c "veronica/common"
)


type Dispatcher struct {
    threads    int
    WorkerPool chan chan int        // use to get aviable workers
    Workers    []*Worker            // use to save active worker instance
}

func NewDispatcher() *Dispatcher {
    pThreads := c.Config.Demo.Threads
    return &Dispatcher{
        threads    : pThreads,
        WorkerPool : make(chan chan int, pThreads),
        Workers    : []*Worker{},
    }
}

// Get dispatcher running status.
func (d *Dispatcher) Status() c.RState {
    for _, worker := range d.Workers {
        if worker.State == c.Running {
            return c.Running
        } else {
            continue
        }
    }
    return c.Stopped
}

// Get dispatcher monitor status.
func (d *Dispatcher) Monitor() []*c.MonitorPack {
    monitorStatus := []*c.MonitorPack{}
    for _, worker := range d.Workers {
        monitorStatus = append(monitorStatus, worker.ProcStatus())
    }
    return monitorStatus
}

// Send Signal to all managed goroutines.
func (d *Dispatcher) Receive(m <-chan c.SIGNAL) {
    go func() {
        for {
            select {
            case signal := <-m:
                if signal == c.SIGSTART {
                    d.run()
                } else if signal == c.SIGSTOP {
                    d.stop()
                }
            default:
                time.Sleep(c.DefaultSleepDur)
            }
        }
    }()
}

// Start running dispatcher.
func (d *Dispatcher) run() {
    for i := 0; i < d.threads; i++ {
        id     := i + 1
        worker := NewWorker(id, d.WorkerPool)
        worker.Start()                            // start worker

        d.Workers = append(d.Workers, worker)     // add worker instance to slice
    }
    go d.dispatch()
}


func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <-c.DemoQueue:
            go func(job int) {
                // Try choose one worker to work, if no aviable worker, it will be blocked.
                workerChan := <-d.WorkerPool
                workerChan<- job        // add this job to this worker Chan
            }(job)
        }
    }
}

func (d *Dispatcher) stop() {
    for _, worker := range d.Workers {
        worker.Stop()
    }
}

type Worker struct {
    Id         int
    State      c.RState
    WorkerPool chan chan int
    JobChan    chan int
}

func NewWorker(id int, workPool chan chan int) *Worker {
    return &Worker{
        Id         : id,
        State      : c.Stopped,
        WorkerPool : workPool,                 // pool chan to save worker job chan
        JobChan    : make(chan int),           // receive task from outside world
    }
}

func (w *Worker) Start() {
    w.State = c.Running
    go w.run()
}

func (w *Worker) run() {
    c.Logger.WithFields(c.LogFields{
        "moduleName" : ModuleName,
        "threadId"   : w.Id,
    }).Info("thread start")
    for w.State == c.Running {
        // regist worker job chan, means the routine is ready to serve
        w.WorkerPool<- w.JobChan

        select {
        case dq := <-w.JobChan:
            c.Logger.WithFields(c.LogFields{
                "moduleName" : ModuleName,
                "threadId"   : w.Id,
            }).Infof("worker received num: %d", dq)
        default:
            time.Sleep(c.DefaultSleepDur)
        }
    }
}

func (w *Worker) Stop() {
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
        "moduleName" : ModuleName,
        "threadId"   : w.Id,
    }).Info("thread stopped")
}

func (w *Worker) ProcStatus() *c.MonitorPack {
    stdLevel  := c.MONITOR_INFO
    if w.State == c.Running {
        stdLevel = c.MONITOR_INFO
    } else {
        stdLevel = c.MONITOR_ERROR
    }

    content := fmt.Sprint(w.State)
    return &c.MonitorPack{
        StdLevel : stdLevel,
        Content  : content,
        Fields   : c.LogFields{
            "threadId" : w.Id,
            "module"   : ModuleName,
            "state"    : w.State,
        },
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
