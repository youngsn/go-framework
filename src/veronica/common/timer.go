package common

import (
    "time"
)

// time cost timer
type TimerBox struct {
    box map[string]*timer
}

type timer struct {
    Stopped bool
    StartTs time.Time
    UsedTs  int64
}

func NewTimerBox() *TimerBox {
    return &TimerBox{
        box : map[string]*timer{},
    }
}

// start a named timer
func (p *TimerBox) Start(name string) bool {
    if _, ok := p.box[name]; !ok { // init timer
        p.box[name] = &timer{
            Stopped : true,
        }
    }

    tm, _ := p.box[name]
    if tm.Stopped == false {       // timer is running
        return false
    }
    tm.Stopped = false
    tm.StartTs = time.Now()
    return true
}

// stop a named timer
func (p *TimerBox) Stop(name string) (int64, bool) {
    tm, ok := p.box[name]
    if !ok {
        return 0, false
    }

    tm.UsedTs  = tm.StartTs.UnixNano() - time.Now().UnixNano()
    tm.Stopped = true
    return tm.UsedTs, true
}

// get timer used time
func (p *TimerBox) GetUseTime(name string) float64 {
    tm, ok := p.box[name]
    if !ok {
        return 0
    }
    usage  := tm.UsedTs
    return float64(usage / 1000000)
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
