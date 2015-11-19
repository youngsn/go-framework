package sys


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
    modulesPriority     []string                     // module priority array
    moduleStopDelay     int                          // module force shutdown waiting time
    pipes               map[string]chan SIGNAL       // module pipe map

    mu                  *sync.Mutex

    SysMonitor          *Monitor                     // monitor module
    SysTimerTask        *TimerTask                   // ticker module
}


func NewModuleManager(m map[string]Module, p []string) *ModuleManager {
    pipes                  := map[string]chan SIGNAL{}
    for moduleName, module := range m {             // init pipe chan, module Ctrl() method get chan pointer
        pipes[moduleName]   = make(chan SIGNAL, 1)
        module.Ctrl(pipes[moduleName])
    }

    return &ModuleManager{
        Modules            : m,
        modulesPriority    : p,
        moduleStopDelay    : 30,                   // waiting 30s then module stop abandoned
        pipes              : pipes,
        mu                 : new(sync.Mutex),
        SysMonitor         : NewMonitor(),
        SysTimerTask       : NewTimerTask(),
    }
}


// Sending broadcast message to all modules.
// Custome SIGNAL can be defined.
func (this *ModuleManager) SendBoardcast(signal SIGNAL) {
    for _, modulePipe := range this.pipes {
        modulePipe<- signal
    }
}


// Send SIGNAL to one module.
func (this *ModuleManager) SendSignal(signal SIGNAL, name string) error {
    if _, ok := this.Modules[name]; !ok {
        return fmt.Errorf("module %s not exist", name)
    }
    this.pipes[name]<- signal

    return nil
}


// Start all modules.
// Method used SendBoardcast() sending START SIGNAL.
// Also here can be changed in sequence as expected.
// Module should start in time, or panic will be throwned.
func (this *ModuleManager) StartModules() {
    this.SendBoardcast(SIGNAL_START)

    // start modules
    started              := false
    for started != true {
        startedModules   := 0
        for _, module := range this.Modules {
            if module.Status() == Running {
                startedModules++
            } else {
                time.Sleep(time.Microsecond * 100)       // if not started, sleep 100ms
            }
        }

        if len(this.Modules) == startedModules {
            started           = true
        }
    }

    this.SysMonitor.Start()
    this.SysTimerTask.Start()
}


// Stop all modules.
// Stop work will follow the defined priority in sequence.
func (this *ModuleManager) StopModules() {
    this.SysMonitor.Stop()
    this.SysTimerTask.Stop()

    for _, moduleName := range this.modulesPriority {
        err           := this.stopModule(moduleName)
        if err != nil {
            Log.Warnf("module %s not stopped, %s", moduleName, err.Error())
        } else {
            Log.Infof("module %s stopped", moduleName)
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

    this.SendSignal(SIGNAL_STOP, name)
    timer       := time.NewTimer(time.Duration(this.moduleStopDelay) * time.Second)

    stopped     := false
    for stopped != true {
        select {
        case <-timer.C:
            this.unloadModule(name)
            stopped            = true
            return fmt.Errorf("can not stop over 30s, check Stop() code carefully")
        default:
            if module.Status() == Stopped {
                this.unloadModule(name)
                stopped        = true
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

    newPriority       := []string{}
    for _, moduleName := range this.modulesPriority {   // delete module from array
        if moduleName != name {
            newPriority     = append(newPriority, moduleName)
        } else {
            continue
        }
    }

    this.modulesPriority    = newPriority
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
