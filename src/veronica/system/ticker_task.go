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
    state    c.RState
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
        state    : c.Stopped,
        interval : 10,
    }
}

// Start TrickerTask instance.
// NOTICE: All task run by goroutine.
func (p *TickerTask) run() {
    p.initTkTask()                      // init timer tickers
    p.state = c.Running
    for _, tk := range p.tkTask {
        go func(task *tkTask) {         // start by goroutine
            for p.state == c.Running {
                select {
                case tm := <-task.Tk.C:
                    p.timer.Start(task.Name)
                    task.TkHd()
                    usage, _ := p.timer.Stop(task.Name)

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
func (p *TickerTask) stop() {
    for _, task := range p.tkTask {
        task.Tk.Stop()
    }
    p.state = c.Stopped
}

// Init ticker tasks from config file.
// Task handler are all from trigger_handler.go.
func (p *TickerTask) initTkTask() {
    for name, interval := range c.Config.Tickers {
        hd   := getTickerHandler(name)               // get task handler
        if interval < p.interval {
            panic(fmt.Sprintf("task: %s, interval less %ds", name, p.interval))
        }

        task := &tkTask{
            Name   : name,
            Status : false,
            Tk     : time.NewTicker(time.Duration(interval) * time.Second),
            TkHd   : hd,
        }
        p.tkTask = append(p.tkTask, task)
        c.Logger.WithFields(c.LogFields{
            "taskName" : name,
            "interval" : interval,
        }).Infof("task: %s, regist success", name)
    }
}

func (p *TickerTask) Init(m <-chan c.SIGNAL) {
    go func() {
        for {
            select {
            case sg := <-m:
                if sg == c.SIGSTART {
                    p.run()
                } else if sg == c.SIGSTOP {
                    p.stop()
                }
            case <- time.After(c.DefaultSleepDur):
            }
        }
    }()
}

func (p *TickerTask) Status() c.RState {
    return p.state
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
