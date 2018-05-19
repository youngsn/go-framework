package system


// TimerTask defines all program timer tasks.
// Timer interval are all defined in config file.
// Handler are all defined in trigger_handler.go.
// Each task will start by using goroutine,
// so tasks will no be effected by others.
// @AUTHOR tangyang
import (
    "fmt"
    "time"
    "strconv"

    c "veronica/common"
)

type TickerTask struct {
    name     string
    timer    *c.TimerBox     // timer
    tkTask   []*tkTask       // task list
    interval int64
    State    c.RState
}

// ticker task struct
type tkTask struct {
    Name   string            // taskname
    Status bool              // run status
    Tk     *time.Ticker      // task ticker interval
    TkHd   TickerHandler     // task handler fuc
}

func NewTickerTask() *TickerTask {
    return &TickerTask{
        name     : "Ticker",
        timer    : c.NewTimerBox(),
        tkTask   : []*tkTask{},
        interval : 10,
        State    : c.Stopped,
    }
}

// Start TrickerTask instance.
// NOTICE: All task run by goroutine.
func (t *TickerTask) run() {
    t.initTkTask()                      // init timer tickers
    t.State = c.Running
    for _, tk := range t.tkTask {
        go func(task *tkTask) {         // start by goroutine
            for t.State == c.Running {
                select {
                case tm := <-task.Tk.C:
                    t.timer.Start(task.Name)
                    task.TkHd()
                    usage, _ := t.timer.Stop(task.Name)

                    c.Logger.WithFields(c.LogFields{
                        "exTime"   : tm.Format("2006-01-02 15:04:05"),
                        "usage"    : strconv.FormatFloat(float64(usage / 1000000), 'f', 2, 32),
                        "taskName" : task.Name,
                    }).Infof("task finished")
                case <-time.After(c.DefaultSleepDur):
                }
            }
        }(tk)
    }
}

// stop ticker tasks.
func (t *TickerTask) stop() {
    for _, task := range t.tkTask {
        task.Tk.Stop()
    }
    t.State = c.Stopped
}

// Init ticker tasks from config file.
// Task handler are all from trigger_handler.go.
func (t *TickerTask) initTkTask() {
    for name, interval := range c.Config.Tickers {
        hd   := getTickerHandler(name)               // get task handler
        if interval < t.interval {
            panic(fmt.Sprintf("%s, interval must large %ds", name, t.interval))
        }

        task := &tkTask{
            Name   : name,
            Status : false,
            Tk     : time.NewTicker(time.Duration(interval) * time.Second),
            TkHd   : hd,
        }
        t.tkTask = append(t.tkTask, task)
        c.Logger.WithFields(c.LogFields{
            "taskName" : name,
            "interval" : interval,
        }).Infof("task: %s, regist success", name)
    }
}

func (t *TickerTask) Init(m <-chan c.SIGNAL) {
    go func() {
        for {
            select {
            case signal := <-m:
                if signal == c.SIGSTART {
                    t.run()
                } else if signal == c.SIGSTOP {
                    t.stop()
                }
            case <- time.After(c.DefaultSleepDur):
            }
        }
    }()
}

func (t *TickerTask) Status() c.RState {
    return t.State
}

func (t *TickerTask) Monitor() []*c.MonitorPack {
    return nil
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
