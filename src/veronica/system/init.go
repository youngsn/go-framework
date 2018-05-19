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
    app.Before  = func(c *cli.Context) error {   // Parse cmd & init yaml config file
        return initConfig(c)
    }
    app.Action  = func(c *cli.Context) {         // Run apps
        err     = appRun()
    }
    app.After   = func(c *cli.Context) error {   // If run failed, handle it easily, just panic
        if err != nil {
            panic(err.Error())
        }
        return err
    }
    return app
}

// Parse cmd params and get yml config filepath
// Then pase config file
func initConfig(context *cli.Context) error {
    cfgFile   := context.String("config")
    if cfgFile == "" {
        return fmt.Errorf("no config file exist, please check filepath")
    }

    var (
        cfg []byte
        err error
        config c.AppConfig
    )
    if cfg, err = ioutil.ReadFile(cfgFile); err != nil {
        return err
    }
    if err = yaml.Unmarshal(cfg, &config); err != nil {
        return err
    }
    c.Config = config
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
    initLogger()

    // start app
    c.Logger.Infof("bootstrap %s, worker %d", c.APP_NAME, c.Config.Worker)
    SysPprof   = NewPprof()             // new pprof monitor
    SysManager = initAppModule()        // init modules

    // init channel
    maxChanSize := c.Config.ChanSize
    c.DemoQueue  = make(chan int, maxChanSize)
    return nil
}

func initAppModule() *ModuleManager {
    m := NewModuleManager()
    m.Init("Demo", demo.NewDispatcher())        // add app modules
    return m
}

func initLogger() {
    flag    := "app"
    c.Logger = c.NewLog(flag)
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
