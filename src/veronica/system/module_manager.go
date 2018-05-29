package system


// ModuleManager manuals all Modules & it's goroutines.
// ModuleManager use one direct chan sending SIGNAL to sub module.
// More Module control methods can be added here, like ReloadModule(), etc.
// @AUTHOR tangyang

import(
    "fmt"
    "time"
    "sync"

    c "veronica/common"
)

// Manager all system & custome modules.
// Go map(search instance)&array(priority defines) are used to Manual.
// One direct chan can send SINGAL to module main goroutines.
type ModuleManager struct {
    mu         *sync.Mutex
    stopDelay  time.Duration               // module force shutdown waiting time
    priority   []string                    // module priority list
    pipes      map[string]chan c.SIGNAL    // module pipe map
    appModules map[string]c.Module         // module map
}

func NewModuleManager() *ModuleManager {
    mg := ModuleManager{
        mu         : new(sync.Mutex),
        stopDelay  : time.Duration(30) * time.Second,
        priority   : []string{},
        pipes      : map[string](chan c.SIGNAL){},
        appModules : map[string]c.Module{},
    }
    // auto regist monitor & ticker task
    mg.Init("Monitor", NewMonitor())
    mg.Init("TickerTask", NewTickerTask())
    return &mg
}

// init module to module manager
func (p *ModuleManager) Init(name string, inst c.Module) {
    p.appModules[name] = inst

    p.priority    = append(p.priority, name)
    p.pipes[name] = make(chan c.SIGNAL, 1)          // init signal chan
    inst.Init(p.pipes[name])
}

// Sending broadcast signal to all modules.
// Custome SIGNAL can be defined.
func (p *ModuleManager) SendBoardcast(sg c.SIGNAL) {
    for _, pipe := range p.pipes {
        pipe<- sg
    }
}

// Send SIGNAL to module
func (p *ModuleManager) SendSignal(sg c.SIGNAL, name string) error {
    if _, ok := p.appModules[name]; !ok {
        return fmt.Errorf("module %s not exist", name)
    }
    p.pipes[name]<- sg
    return nil
}

// Start all modules.
// Method used SendBoardcast() sending START SIGNAL.
// Also here can be changed in sequence as expected.
// Module should start in time, or panic will be throwned.
func (p *ModuleManager) Start() error {
    if len(p.appModules) == 0 {
        return fmt.Errorf("failed, no valid modules")
    }
    for name, inst := range p.appModules {
        p.SendSignal(c.SIGSTART, name)
        for inst.Status() != c.Running {            // wait module to startup
            time.Sleep(c.DefaultSleepDur)
        }

        c.Logger.WithFields(c.LogFields{
            "module" : name,
        }).Infof("%s, ready to work", name)
    }
    return nil
}

// Stop all modules.
// Stop work will follow the defined priority in sequence.
func (p *ModuleManager) Stop() {
    pri := make([]string, len(p.priority))
    copy(pri, p.priority)
    for _, name := range pri {
        if err  := p.stopModule(name); err != nil {
            c.Logger.WithFields(c.LogFields{
                "module" : name,
                "errmsg" : err.Error(),
            }).Errorf("%s, stop failed", name)
        } else {
            c.Logger.WithFields(c.LogFields{
                "module" : name,
            }).Infof("%s, stopped", name)
        }
    }
}

// Get app module list
func (p *ModuleManager) GetAppModules() map[string]c.Module {
    return p.appModules
}

// Get Started module quantity.
func (p *ModuleManager) runModuleNum() int {
    cnt := 0
    for _, module := range p.appModules {
        if module.Status() == c.Running {
            cnt++
        }
    }
    return cnt
}

// Get Stopped module qunatity.
func (p *ModuleManager) stopModuleNum() int {
    cnt := 0
    for _, inst := range p.appModules {
        if inst.Status() == c.Stopped {
            cnt++
        }
    }
    return cnt
}

// Stop one module.
// Send STOP SIGNAL to module and wait module stop.
// If stopping used more than this.moduleStopDelay time,
// ModuleManager will forget this module status and flag it stopped.
func (p *ModuleManager) stopModule(name string) error {
    inst, ok := p.appModules[name]
    if ok == false {
        return fmt.Errorf("%s not exist", name)
    }

    p.SendSignal(c.SIGSTOP, name)
    timer := time.NewTimer(p.stopDelay)
    stop  := false
    for !stop {
        select {
        case <-timer.C:
            p.unload(name)
            stop = true
            return fmt.Errorf("can not stop in 30s, abandon")
        default:
            if inst.Status() == c.Stopped {
                p.unload(name)
                stop = true
            } else {
                time.Sleep(c.DefaultSleepDur)
            }
        }
    }
    return nil
}

// Unload module, used to stop module.(just delete module from map and priority)
func (p *ModuleManager) unload(name string) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if _, ok := p.appModules[name]; ok {
        delete(p.appModules, name)
    }

    for i, moduleName := range p.priority {     // delete module from array
        if moduleName == name {                 // delete
            p.priority = append(p.priority[: i], p.priority[(i + 1) : ]...)
            break
        } else {
            continue
        }
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
