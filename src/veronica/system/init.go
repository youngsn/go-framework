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

// In init func, use cli package to parse config and manage app life cycle
func Init() *cli.App {
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
    app.Before  = func(ct *cli.Context) error {   // init config from -c & init app
        if err = initConfig(ct); err != nil {
            return err
        }
        if err = initApp(); err != nil {
            return err
        }
        return nil
    }
    app.Action  = func(ct *cli.Context) {         // run app
        c.Logger.Infof("bootstrap %s, worker %d", c.APP_NAME, c.Config.Worker)
        err     = runApp()
    }
    app.After   = func(ct *cli.Context) error {   // if app run failed, handle it, panic to user
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

// App run entry, start from here
func runApp() error {
    if err := c.WritePid(); err != nil {
        return err
    }
    worker := c.Config.Worker       // work cpu num
    if worker > runtime.NumCPU() {
        worker = runtime.NumCPU()
    }
    runtime.GOMAXPROCS(worker)

    SysPprof.WebMonitor()                         // pprof monitor
    if err := SysManager.Start(); err != nil {    // start app modules
        return err
    }
    c.Logger.Infof("app start success")

    NewSignal().Start()             // signal capture, no end loop
    c.UnlinkPid()
    return nil
}

// Init app data.
func initApp() error {
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
    SysPprof   = NewPprof()             // pprof monitor
    SysManager = initAppModule()        // init app module

    chanSize  := c.Config.ChanSize
    c.MonitorQueue = make(chan []*c.MonitorPack, chanSize)
    c.DemoQueue    = make(chan int, chanSize)
    return nil
}

func initAppModule() *ModuleManager {
    m := NewModuleManager()
    m.Init("Demo", demo.NewDemoDispatcher())        // add app modules
    return m
}

func initLogger() {
    flag    := "app"
    c.Logger = c.NewLog(flag)
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
