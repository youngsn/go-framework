package system


// ModuleManager manuals all Modules & it's goroutines.
// ModuleManager use one direct chan sending SIGNAL to sub module.
// More Module control methods can be added here, like ReloadModule(), etc.
// @AUTHOR tangyang

import(
    "fmt"
    "time"
    "sync"

    . "veronica/common"

    Log "github.com/cihub/seelog"
)


// Manager all system & custome modules.
// Go map(search instance)&array(priority defines) are used to Manual.
// One direct chan can send SINGAL to module main goroutines.
type ModuleManager struct {
    Modules             map[string]Module            // custome module map
    priority            []string                     // module priority array
    moduleStopDelay     int                          // module force shutdown waiting time
    pipes               map[string]chan SIGNAL       // module pipe map

    mu                  *sync.Mutex

    SysMonitor          *Monitor                     // monitor module
    SysTimerTask        *TimerTask                   // ticker module
}


func NewModuleManager() *ModuleManager {
    return &ModuleManager{
        Modules            : map[string]Module{},
        priority           : []string{},
        moduleStopDelay    : 30,                   // waiting 30s then module stop abandoned
        pipes              : map[string]chan SIGNAL{},
        mu                 : new(sync.Mutex),
        SysMonitor         : NewMonitor(),
        SysTimerTask       : NewTimerTask(),
    }
}


// init module to module manager
func (this *ModuleManager) InitModule(name string, m Module) {
    this.Modules[name]     = m
    this.priority          = append(this.priority, name)

    // init pipe chan, module Ctrl() method get chan pointer
    this.pipes[name]       = make(chan SIGNAL, 1)
    m.Ctrl(this.pipes[name])
}


// Sending broadcast message to all modules.
// Custome SIGNAL can be defined.
func (this *ModuleManager) SendBoardcast(s SIGNAL) {
    for _, modulePipe := range this.pipes {
        modulePipe<- s
    }
}


// Send SIGNAL to one module.
func (this *ModuleManager) SendSignal(s SIGNAL, name string) error {
    if _, ok := this.Modules[name]; !ok {
        return fmt.Errorf("module %s not exist", name)
    }
    this.pipes[name]<- s

    return nil
}


// Start all modules.
// Method used SendBoardcast() sending START SIGNAL.
// Also here can be changed in sequence as expected.
// Module should start in time, or panic will be throwned.
func (this *ModuleManager) StartModules() error {
    if len(this.Modules) == 0 {
        return fmt.Errorf("Start failed, no module exist")
    }

    this.SendBoardcast(SIGSTART)

    // start modules
    started           := false
    for started != true {
        startModule   := 0
        for _, module := range this.Modules {
            if module.Status() == Running {
                startModule++
            } else {
                time.Sleep(time.Microsecond * 100)       // if not started, sleep 100ms
            }
        }

        if len(this.Modules) == startModule {
            started           = true
        }
    }

    this.SysMonitor.Start()
    this.SysTimerTask.Start()

    return nil
}


// Stop all modules.
// Stop work will follow the defined priority in sequence.
func (this *ModuleManager) StopModules() {
    this.SysMonitor.Stop()
    this.SysTimerTask.Stop()

    for _, name := range this.priority {
        err     := this.stopModule(name)
        if err  != nil {
            Log.Warnf("module %s not stopped, %s", name, err.Error())
        } else {
            Log.Infof("module %s stopped", name)
        }
    }
}


// Stop one module.
// Send STOP SIGNAL to module and wait module stop.
// If stopping used more than this.moduleStopDelay time,
// ModuleManager will forget this module status and flag it stopped.
func (this *ModuleManager) stopModule(name string) error {
    module, ok  := this.Modules[name]
    if ok == false {
        return fmt.Errorf("%s not exist", name)
    }

    this.SendSignal(SIGSTOP, name)
    timer       := time.NewTimer(time.Duration(this.moduleStopDelay) * time.Second)

    stopped     := false
    for stopped != true {
        select {
        case <-timer.C:
            this.unloadModule(name)
            stopped             = true
            return fmt.Errorf("can not stop over 30s, check Stop() code carefully")
        default:
            if module.Status() == Stopped {
                this.unloadModule(name)
                stopped         = true
            } else {
                time.Sleep(time.Microsecond * 100)
            }
        }
    }

    return nil
}


// Unload module, used to stop module.(just delete module from map and priority)
func (this *ModuleManager) unloadModule(name string) {
    this.mu.Lock()
    defer this.mu.Unlock()

    if _, ok := this.Modules[name]; ok {
        delete(this.Modules, name)
    }

    for i, moduleName := range this.priority {   // delete module from array
        if moduleName == name {                  // delete
            this.priority   = append(this.priority[:i], this.priority[i+1:]...)
            break
        } else {
            continue
        }
    }
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
