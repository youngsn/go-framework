package common


import (
    "time"
)


// This file can added all the constants that all package can use them.
// Just import the package as . "common" when you want to use them.
// Maybe package are used in every other packages.

// system environment defines
const APP_NAME    = "veronica"
const APP_VERSION = "2.1.5"

var (
    Environ   string
    StartTime time.Time
    SysPath   string                // app base path
    RunPath   string                // app run path
    DataPath  string                // app data path
    Config    AppConfig             // app global config
    Logger    *Log                  // log engine
    NameCli   *NameClient           // naming client
)

// Project channels, THIS IS JUST DEMO.
var (
    DemoQueue chan int
    MonitorQueue chan []*MonitorPack
)

// All constant defines are below.
// System consts defines.
const (
    DefaultSleepDur time.Duration = 100 * time.Microsecond      // select default sleep time
)

// System module signal defines.
// You can add more own SIGNALs for modules.
type SIGNAL int
const (
    SIGSTART SIGNAL = iota + 1
    SIGSTOP
    SIGRELOAD
    SIGMONITOR
)

func (s SIGNAL) String() string {
    switch s {
    case SIGSTART:
        return "Start"
    case SIGSTOP:
        return "Stop"
    case SIGRELOAD:
        return "Reload"
    case SIGMONITOR:
        return "Monitor"
    default:
        return "unknown"
    }
}

// Module running status.
type RState int
const (
    Running RState = iota + 1
    Stopped
    Waiting
)

func (s RState) String() string {
    switch s {
    case Running:
        return "Running"
    case Stopped:
        return "Stopped"
    case Waiting:
        return "Waiting"
    default:
        return "Unknown"
    }
}

// Modules monitor log level defines.
const (
    MONITOR_INFO int = iota + 1
    MONITOR_ERROR
)

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
