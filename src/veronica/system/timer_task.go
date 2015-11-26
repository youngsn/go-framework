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

    . "veronica/common"

    Log "github.com/cihub/seelog"
)


// One TickerTask struct
type TickerTask struct {
    Name          string            // taskname
    Status        bool              // running status
    Tk            *time.Ticker      // task ticker
    TickerHd      TickerHandler     // task handler fuc
}


type TimerTask struct {
    State         RState

    tickerTask    []*TickerTask             // task list
    minInterval   int64
}


func NewTimerTask() *TimerTask {
    return &TimerTask{
        tickerTask   : []*TickerTask{},
        State        : Stopped,
        minInterval  : 10,
    }
}


// Start TriggerTasks instance.
// NOTICE: All task are started by goroutine.
func (this *TimerTask) Start() {
    this.initTickerTask()                   // init timer tickers

    this.State   = Running
    for _, task := range this.tickerTask {
        go func(task *TickerTask) {         // start by goroutine
            for this.State == Running {
                select {
                case t := <-task.Tk.C:
                    Log.Infof("Ticker task %s @(%s)", task.Name, t.Format("2006-01-02 15:04:05"))

                    s              := time.Now().UnixNano()
                    task.TickerHd()
                    e              := time.Now().UnixNano()

                    usage          := strconv.FormatFloat(float64((e - s) / 1000000), 'f', 2, 32)
                    Log.Infof("Task %s finished, time: %s ms", task.Name, usage)
                default:
                    time.Sleep(DefaultSleepDur)
                }
            }
        }(task)
    }

    Log.Infof("Tricker Thread, start")
}


// Init ticker tasks from config file.
// Task handler are all from trigger_handler.go.
func (this *TimerTask) initTickerTask() {
    for name, ticker := range Config.Tickers {
        handler      := getTickerHandler(name)          // get task handler
        if ticker.Interval < this.minInterval {
            panic(fmt.Sprintf("Task %s, interval must large %ds", name, this.minInterval))
        }
        tk           := time.NewTicker(time.Second * time.Duration(ticker.Interval))

        t            := &TickerTask{
            Name      : name,
            Status    : false,
            Tk        : tk,
            TickerHd  : handler,
        }

        this.tickerTask     = append(this.tickerTask, t)
        Log.Infof("Regist ticker task: %s, interval: %ds", name, ticker.Interval)
    }
}


// Stop all tasks.
func (this *TimerTask) Stop() {
    this.State    = Stopped
    for _,  task := range this.tickerTask {
        task.Tk.Stop()
    }

    Log.Infof("Tricker thread, stopped")
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
