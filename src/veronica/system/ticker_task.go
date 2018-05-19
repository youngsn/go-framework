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

// TickerTask struct
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
func (t *TickerTask) Start() {
    t.initTkTask()                      // init timer tickers
    t.State = c.Running
    for _, task := range t.tkTask {
        go func(task *tkTask) {         // start by goroutine
            for t.State == c.Running {
                select {
                case tm := <-task.Tk.C:
                    t.timer.Start(task.Name)
                    task.TkHd()
                    usage, _ := t.timer.Stop(task.Name)

                    c.Logger.WithFields(c.LogFields{
                        "exTime" : tm.Format("2006-01-02 15:04:05"),
                        "usage"  : strconv.FormatFloat(float64(usage / 1000000), 'f', 2, 32),
                        "taskName" : task.Name,
                    }).Infof("task finished")
                case <-time.After(c.DefaultSleepDur):
                }
            }
        }(task)
    }
    c.Logger.Infof("%s, start work", t.name)
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

// Stop all tasks.
func (t *TickerTask) Stop() {
    t.State = c.Stopped
    for _, task := range t.tkTask {
        task.Tk.Stop()
    }
    c.Logger.Infof("%s, stopped", t.name)
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
