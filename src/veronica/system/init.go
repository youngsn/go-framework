package system


// Program & global params are inits here.
import (
    "os"
    "fmt"
    "runtime"
    "io/ioutil"

    "veronica/demo"
    c "veronica/common"

    "gopkg.in/yaml.v2"
    "github.com/codegangsta/cli"
    Log "github.com/cihub/seelog"
)


// system control global vars
var (
    SysPprofMonitor *PprofMonitor       // pprof var
    SysManager      *ModuleManager      // module manage
)

func InitApp() *cli.App {
    cmdFlag := []cli.Flag{
        cli.StringFlag{
            Name  : "config, c",
            Value : "",
            Usage : "app yaml config file",
        },
    }

    var err error
    app        := cli.NewApp()
    app.Name    = c.APP_NAME
    app.Version = c.APP_VERSION
    app.Usage   = fmt.Sprintf("Run %s service", c.APP_NAME)
    app.Flags   = cmdFlag
    app.Before  = func(c *cli.Context) error {   // Parse cmd params
        return parseParams(c)
    }
    app.After   = func(c *cli.Context) error {   // If has exec errors, return
        if err != nil {             // run failed, we need panic
            panic(err.Error())
        }
        return err
    }
    app.Action  = func(c *cli.Context) {         // Run apps 
        err     = appRun()
    }
    return app
}

// Parse cmd params & do param handle.
func parseParams(context *cli.Context) error {
    cfgFile    := context.String("config")
    if cfgFile == "" {
        return fmt.Errorf("no config file found, please specify")
    }

    config    := c.ConfigStruct{}
    if c, err := ioutil.ReadFile(cfgFile); err != nil {
        return err
    } else {
        if err = yaml.Unmarshal(c, &config); err != nil {
            return err
        }
    }
    c.Config   = config
    return nil
}

// app run main entry.
func appRun() error {
    if err := appInit(); err != nil {
        return err
    }
    if err := c.WritePid(); err != nil {
        return err
    }

    processors := c.Config.Global.Processors
    runtime.GOMAXPROCS(processors)

    SysPprofMonitor.WebPprofMonitor()                  // pprof monitor
    if err := SysManager.StartModules(); err != nil {  // start modules
        return err
    }

    sysSignal := NewSignal()                           // system signal capture
    sysSignal.Start()

    c.UnlinkPid()
    return nil
}

// Init app data. 
func appInit() error {
    // init common package global vars
    sysPath, _   := os.Getwd()
    c.Environment = c.Config.Environment
    c.SysPath     = sysPath + "/"
    c.RunPath     = sysPath + "/run/"
    if err := c.FilePathExist(c.RunPath); err != nil {
        return err
    }
    if err := loggerInit(); err != nil {            // init log engine
        return err
    }

    Log.Infof("Bootstrap %s", c.APP_NAME)           // Program start here
    SysPprofMonitor    = NewPprofMonitor()          // new pprof monitor
    SysManager         = modulesInit()              // new module & module manager init

    // chan init
    maxSize           := c.Config.Global.MaxChannelSize
    c.DemoQueue        = make(chan int, maxSize)

    return nil
}

// Init custome modules here.
func modulesInit() *ModuleManager {
    m := NewModuleManager()
    m.InitModule("Demo", demo.NewDManager())
    return m
}

// Init log engine instance.
// All defined loggers are init into LoggerFactory at once.
// Maybe goroutine used different instances.
// Log "github.com/cihub/seelog" used default config.
func loggerInit() error {
    if err := initLoggerFactory(); err != nil {
        return err
    }

    Log.ReplaceLogger(c.GetLogger("default"))       // default logger
    demo.Log = c.GetLogger("demo")

    return nil
}

// Init logger factory.
func initLoggerFactory() error {
    c.LoggerFactory         = map[string]Log.LoggerInterface{}
    for loggerName, config := range c.Config.Logger {      // convert logger config to inst
        if logInstance, err := Log.LoggerFromConfigAsString(config); err == nil {
            c.LoggerFactory[loggerName] = logInstance
        } else {
            return err
        }
    }

    return nil
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
