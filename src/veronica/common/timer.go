package common

import(
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

/**
 * start a named timer
 *
 * @param name
 * @return bool
 */
func (t *TimerBox) Start(name string) bool {
    if _, ok := t.box[name]; !ok { // init timer
        t.box[name] = &timer{
            Stopped : true,
        }
    }

    tm, _ := t.box[name]
    if tm.Stopped == false {       // timer is running
        return false
    }
    tm.Stopped = false
    tm.StartTs = time.Now()
    return true
}

/**
 * stop a timer
 *
 * @param name
 * @return int64, bool
 */
func (t *TimerBox) Stop(name string) (int64, bool) {
    tm, ok := t.box[name]
    if !ok {                                // init timer
        return 0, false
    }

    tm.UsedTs  = tm.StartTs.UnixNano() - time.Now().UnixNano()
    tm.Stopped = true
    return tm.UsedTs, true
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
