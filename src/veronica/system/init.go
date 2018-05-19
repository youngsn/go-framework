package system


// app & global vals init here.
import (
    "os"
    "fmt"
    "runtime"
    "io/ioutil"

    "veronica/demo"
    c "veronica/common"

    "gopkg.in/yaml.v2"
    "github.com/codegangsta/cli"
)

// system control global vars
var (
    SysPprof   *PprofMonitor       // pprof monitor
    SysManager *ModuleManager      // module manage
)

func InitApp() *cli.App {
    cmdFlag := []cli.Flag{
        cli.StringFlag{
            Name  : "config, c",
            Value : "",
            Usage : "app yaml config filepath",
        },
    }

    var err error
    app        := cli.NewApp()
    app.Name    = c.APP_NAME
    app.Version = c.APP_VERSION
    app.Usage   = fmt.Sprintf("Run app %s", c.APP_NAME)
    app.Flags   = cmdFlag
    app.Before  = func(c *cli.Context) error {   // Parse cmd params
        return parseConfig(c)
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

// Parse cmd params and get yml config filepath
// Then pase config file
func parseConfig(context *cli.Context) error {
    cfgFile   := context.String("config")
    if cfgFile == "" {
        return fmt.Errorf("no config file exist, please check filepath")
    }

    var config c.AppConfig
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

    worker := c.Config.Worker       // worker num
    runtime.GOMAXPROCS(worker)
    SysPprof.WebMonitor()           // pprof monitor
    if err := SysManager.Start(); err != nil {     // start app modules
        return err
    }
    c.Logger.Infof("app start success")
    NewSignal().Start()             // signal capture
    c.UnlinkPid()
    return nil
}

// Init app data.
func appInit() error {
    // init common package global vars
    sysPath, _ := os.Getwd()
    c.Environ   = c.Config.Environ
    c.SysPath   = sysPath
    c.RunPath   = sysPath + "/var"
    c.DataPath  = sysPath + "/data"
    if err := c.FilePathExist(c.RunPath); err != nil {
        return err
    }
    if err := c.FilePathExist(c.DataPath); err != nil {
        return err
    }
    // init log engine
    if err := loggerInit(); err != nil {
        return err
    }

    // start app
    c.Logger.Infof("bootstrap %s, worker %d", c.APP_NAME, c.Config.Worker)
    SysPprof   = NewPprof()             // new pprof monitor
    SysManager = initModules()          // new module & module manager init

    // init channel
    maxChanSize := c.Config.ChanSize
    c.DemoQueue  = make(chan int, maxChanSize)
    return nil
}

// Init custome modules here.
func initModules() *ModuleManager {
    m := NewModuleManager()
    m.Init("Demo", demo.NewDispatcher())
    return m
}

// Init global log engine instance.
func loggerInit() error {
    flag    := "app"
    c.Logger = c.NewLog(flag)
    return nil
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
