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
    stopDelay  int                         // module force shutdown waiting time
    priority   []string                    // module priority list
    pipes      map[string]chan c.SIGNAL    // module pipe map
    appModules map[string]c.Module         // module map
}

func NewModuleManager() *ModuleManager {
    mg := ModuleManager{
        mu         : new(sync.Mutex),
        stopDelay  : 30,                            // stop abandon time delay
        priority   : []string{},
        pipes      : map[string](chan c.SIGNAL){},
        appModules : map[string]c.Module{},
    }
    mg.Init("Monitor", NewMonitor())
    mg.Init("TickerTask", NewTickerTask())
    return &mg
}

// init module to module manager
func (m *ModuleManager) Init(name string, inst c.Module) {
    m.appModules[name] = inst

    m.priority    = append(m.priority, name)
    m.pipes[name] = make(chan c.SIGNAL, 1)          // init signal chan
    inst.Init(m.pipes[name])
}

// Sending broadcast signal to all modules.
// Custome SIGNAL can be defined.
func (m *ModuleManager) SendBoardcast(s c.SIGNAL) {
    for _, pipe := range m.pipes {
        pipe<- s
    }
}

// Send SIGNAL to module
func (m *ModuleManager) SendSignal(s c.SIGNAL, name string) error {
    if _, ok := m.appModules[name]; !ok {
        return fmt.Errorf("module %s not exist", name)
    }
    m.pipes[name]<- s
    return nil
}

// Start all modules.
// Method used SendBoardcast() sending START SIGNAL.
// Also here can be changed in sequence as expected.
// Module should start in time, or panic will be throwned.
func (m *ModuleManager) Start() error {
    if len(m.appModules) == 0 {
        return fmt.Errorf("failed, no valid modules")
    }
    for name, inst := range m.appModules {
        m.SendSignal(c.SIGSTART, name)
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
func (m *ModuleManager) Stop() {
    p := make([]string, len(m.priority))
    copy(p, m.priority)
    for _, name := range p {
        if err  := m.stopModule(name); err != nil {
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
func (m *ModuleManager) GetAppModules() map[string]c.Module {
    return m.appModules
}

// Get Started module quantity.
func (m *ModuleManager) startModuleNum() int {
    cnt := 0
    for _, module := range m.appModules {
        if module.Status() == c.Running {
            cnt++
        }
    }
    return cnt
}

// Get Stopped module qunatity.
func (m *ModuleManager) StoppedModuleNum() int {
    cnt := 0
    for _, inst := range m.appModules {
        if inst.Status() == c.Stopped {
            cnt++
        }
    }
    return cnt
}

// Get Manager manage module quantity.
func (m *ModuleManager) ModuleNum() int {
    return len(m.appModules)
}

// Stop one module.
// Send STOP SIGNAL to module and wait module stop.
// If stopping used more than this.moduleStopDelay time,
// ModuleManager will forget this module status and flag it stopped.
func (m *ModuleManager) stopModule(name string) error {
    inst, ok := m.appModules[name]
    if ok == false {
        return fmt.Errorf("%s not exist", name)
    }

    m.SendSignal(c.SIGSTOP, name)
    timer   := time.NewTimer(time.Duration(m.stopDelay) * time.Second)
    stopped := false
    for stopped != true {
        select {
        case <-timer.C:
            m.unload(name)
            stopped = true
            return fmt.Errorf("can not stop in 30s, abandon")
        default:
            if inst.Status() == c.Stopped {
                m.unload(name)
                stopped = true
            } else {
                time.Sleep(c.DefaultSleepDur)
            }
        }
    }
    return nil
}

// Unload module, used to stop module.(just delete module from map and priority)
func (m *ModuleManager) unload(name string) {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, ok := m.appModules[name]; ok {
        delete(m.appModules, name)
    }

    for i, moduleName := range m.priority {     // delete module from array
        if moduleName == name {                 // delete
            m.priority = append(m.priority[: i], m.priority[(i + 1) : ]...)
            break
        } else {
            continue
        }
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
