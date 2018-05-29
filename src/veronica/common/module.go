package common

import (
)

// General dispatcher define.
// It is a program model to module, use to manange multi threads instance
type Dispatcher struct {
    Threads int
    Workers []WorkerInst
}

func NewDispatcher(threads int, workers []WorkerInst) Dispatcher {
    return Dispatcher{
        Threads : threads,
        Workers : workers,
    }
}

// Get dispatcher running status.
func (p *Dispatcher) Status() RState {
    for _, inst := range p.Workers {
        if inst.RunStatus() == Running {
            return Running
        } else {
            continue
        }
    }
    return Stopped
}

// Get worker monitor status
func (p *Dispatcher) Monitor() {
    pk := []*MonitorPack{}
    for _, inst := range p.Workers {
        pk = append(pk, inst.Monitor())
    }
    MonitorQueue<- pk
}


type Worker struct {
    Id    int
    Name  string
    State RState
}

func NewWorker(id int, name string) Worker {
    return Worker{
        Id    : id,
        Name  : name,
        State : Stopped,
    }
}

func (p *Worker) RunStatus() RState {
    return p.State
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
