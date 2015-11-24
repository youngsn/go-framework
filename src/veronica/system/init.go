package system


// Program & global param inits here.
import (
    "os"
    "flag"

    "veronica/demo"
    c "veronica/common"

    "github.com/BurntSushi/toml"
    Log "github.com/cihub/seelog"
)


// system control global vars
var (
    SysPprofMonitor     *PprofMonitor       // pprof var
    SysManager          *ModuleManager      // module manage
)


// Config file parse & global vars init(in common & there)
func Initialize() error {
    if err := parseConfig(); err != nil {   // parse config files
        return err
    }

    // init common package global vars
    c.Name             = c.Config.Name
    c.Environment      = c.Config.Environment
    c.SysPath, _       = os.Getwd()
    c.RunPath          = c.SysPath + "/run/"
    if err := c.FileExist(c.RunPath); err != nil {
        return err
    }

    loggerInit()                                    // init log engine
    Log.Infof("Bootstrap %s", c.Name)               // Program start here
    SysPprofMonitor    = NewPprofMonitor()          // new pprof monitor
    SysManager         = modulesInit()              // new module & module manager init

    // public chan init
    maxSize           := c.Config.Global.MaxChannelSize
    c.DemoQueue        = make(chan int, maxSize)

    return nil
}


// Init custome modules here.
func modulesInit() *ModuleManager {
    modules                 := map[string]c.Module{}
    priority                := []string{}

    modules["demo"]          = demo.NewDManager()
    priority                 = append(priority, "demo")

    return NewModuleManager(modules, priority)
}


// Parse & init Config var to program.
func parseConfig() error {
    var cfgFile         = flag.String("c", "conf.toml", "optimus config file")
    flag.Parse()

    if _, err := toml.DecodeFile(*cfgFile, &c.Config); err != nil {
        return err
    }

    return nil
}


// Init log engine factory.
// All defined loggers are init into LoggerFactory at once.
// Maybe goroutine used different instances.
// Log "github.com/cihub/seelog" used default config.
func loggerInit() {
    c.LoggerFactory          = map[string]Log.LoggerInterface{}
    for loggerName, config  := range c.Config.Logger {      // convert logger config to inst
        if logInstance, err := Log.LoggerFromConfigAsString(config.Conf); err == nil {
            c.LoggerFactory[loggerName] = logInstance
        } else {
            panic(err.Error())
        }
    }

    Log.ReplaceLogger(c.GetLogger("default"))
    demo.Log                 = c.GetLogger("demo")
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
