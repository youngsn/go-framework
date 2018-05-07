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


// One TickerTask struct
type TickerTask struct {
    Name     string            // taskname
    Status   bool              // running status
    Tk       *time.Ticker      // task ticker
    TickerHd TickerHandler     // task handler fuc
}


type TimerTask struct {
    State       c.RState
    tickerTask  []*TickerTask             // task list
    minInterval int64
}

func NewTimerTask() *TimerTask {
    return &TimerTask{
        tickerTask  : []*TickerTask{},
        State       : c.Stopped,
        minInterval : 10,
    }
}

// Start TriggerTasks instance.
// NOTICE: All task are started by goroutine.
func (t *TimerTask) Start() {
    t.initTickerTask()                   // init timer tickers

    t.State = c.Running
    for _, task := range t.tickerTask {
        go func(task *TickerTask) {         // start by goroutine
            for t.State == c.Running {
                select {
                case t := <-task.Tk.C:
                    s  := time.Now().UnixNano()
                    task.TickerHd()
                    e  := time.Now().UnixNano()
                    usage := strconv.FormatFloat(float64((e - s) / 1000000), 'f', 2, 32)

                    c.Logger.WithFields(c.LogFields{
                        "exTime" : t.Format("2006-01-02 15:04:05"),
                        "usage"  : usage,
                        "taskName" : task.Name,
                    }).Info("task exec success")
                default:
                    time.Sleep(c.DefaultSleepDur)
                }
            }
        }(task)
    }
    c.Logger.Info("tricker, start success")
}

// Init ticker tasks from config file.
// Task handler are all from trigger_handler.go.
func (t *TimerTask) initTickerTask() {
    for name, interval := range c.Config.Tickers {
        handler := getTickerHandler(name)               // get task handler
        if interval < t.minInterval {
            panic(fmt.Sprintf("%s, interval must large %ds", name, t.minInterval))
        }

        task := &TickerTask{
            Name     : name,
            Status   : false,
            Tk       : time.NewTicker(time.Second * time.Duration(interval)),
            TickerHd : handler,
        }
        t.tickerTask = append(t.tickerTask, task)
        c.Logger.WithFields(c.LogFields{
            "taskName" : name,
            "interval" : interval,
        }).Info("regist task success")
    }
}

// Stop all tasks.
func (t *TimerTask) Stop() {
    t.State = c.Stopped
    for _,  task := range t.tickerTask {
        task.Tk.Stop()
    }
    c.Logger.Info("tricker, stop success")
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
