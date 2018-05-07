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
    priority  []string                      // module priority array
    stopDelay int                           // module force shutdown waiting time
    pipes     map[string]chan c.SIGNAL      // module pipe map
    mu        *sync.Mutex

    Modules      map[string]c.Module        // run modules
    SysMonitor   *Monitor                   // monitor module
    SysTimerTask *TimerTask                 // ticker module
}

func NewModuleManager() *ModuleManager {
    return &ModuleManager{
        Modules   : map[string]c.Module{},
        priority  : []string{},
        stopDelay : 30,                     // waiting 30s then module stop abandoned
        pipes     : map[string]chan c.SIGNAL{},
        mu        : new(sync.Mutex),
        SysMonitor   : NewMonitor(),
        SysTimerTask : NewTimerTask(),
    }
}

// init module to module manager
func (m *ModuleManager) InitModule(name string, module c.Module) {
    m.Modules[name] = module
    m.priority      = append(m.priority, name)

    // init pipe chan, module Receive() method get chan pointer
    m.pipes[name]   = make(chan c.SIGNAL, 1)
    module.Receive(m.pipes[name])
}

// Sending broadcast message to all modules.
// Custome SIGNAL can be defined.
func (m *ModuleManager) SendBoardcast(s c.SIGNAL) {
    for _, modulePipe := range m.pipes {
        modulePipe<- s
    }
}

// Send SIGNAL to one module.
func (m *ModuleManager) SendSignal(s c.SIGNAL, name string) error {
    if _, ok := m.Modules[name]; !ok {
        return fmt.Errorf("module %s not exist", name)
    }
    m.pipes[name]<- s
    return nil
}

// Start all modules.
// Method used SendBoardcast() sending START SIGNAL.
// Also here can be changed in sequence as expected.
// Module should start in time, or panic will be throwned.
func (m *ModuleManager) StartModules() error {
    if len(m.Modules) == 0 {
        return fmt.Errorf("start failed, module not exist")
    }
    m.SendBoardcast(c.SIGSTART)

    // start module
    started := false
    for started != true {
        if len(m.Modules) == m.StartedModuleNum() {
            started = true
        } else {
            time.Sleep(c.DefaultSleepDur)
        }
    }
    m.SysMonitor.Start()
    m.SysTimerTask.Start()
    return nil
}

// Stop all modules.
// Stop work will follow the defined priority in sequence.
func (m *ModuleManager) StopModules() {
    m.SysMonitor.Stop()
    m.SysTimerTask.Stop()

    p := make([]string, len(m.priority))
    copy(p, m.priority)
    for _, name := range p {
        if err  := m.stopModule(name); err != nil {
            c.Logger.WithFields(c.LogFields{
                "module" : name,
                "err"    : err.Error(),
            }).Error("stop failed")
        } else {
            c.Logger.WithFields(c.LogFields{
                "module" : name,
            }).Info("stop success")
        }
    }
}

// Get Started module quantity.
func (m *ModuleManager) StartedModuleNum() int {
    cnt := 0
    for _, module := range m.Modules {
        if module.Status() == c.Running {
            cnt++
        }
    }
    return cnt
}

// Get Stopped module qunatity.
func (m *ModuleManager) StoppedModuleNum() int {
    cnt := 0
    for _, module := range m.Modules {
        if module.Status() == c.Stopped {
            cnt++
        }
    }
    return cnt
}

// Get Manager manage module quantity.
func (m *ModuleManager) ModulsNum() int {
    return len(m.Modules)
}

// Stop one module.
// Send STOP SIGNAL to module and wait module stop.
// If stopping used more than this.moduleStopDelay time,
// ModuleManager will forget this module status and flag it stopped.
func (m *ModuleManager) stopModule(name string) error {
    module, ok := m.Modules[name]
    if ok == false {
        return fmt.Errorf("%s not exist", name)
    }

    m.SendSignal(c.SIGSTOP, name)
    timer   := time.NewTimer(time.Duration(m.stopDelay) * time.Second)
    stopped := false
    for stopped != true {
        select {
        case <-timer.C:
            m.deleteModule(name)
            stopped = true
            return fmt.Errorf("can not stop in 30s, abandon")
        default:
            if module.Status() == c.Stopped {
                m.deleteModule(name)
                stopped = true
            } else {
                time.Sleep(time.Microsecond * 100)
            }
        }
    }
    return nil
}

// Unload module, used to stop module.(just delete module from map and priority)
func (m *ModuleManager) deleteModule(name string) {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, ok := m.Modules[name]; ok {
        delete(m.Modules, name)
    }

    for i, moduleName := range m.priority {     // delete module from array
        if moduleName == name {                 // delete
            m.priority = append(m.priority[:i], m.priority[i+1:]...)
            break
        } else {
            continue
        }
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
