package common


import (
    "time"

    "github.com/cihub/seelog"
)


// This file can added all the constants that all package can use them.
// Just import the package as . "common" when you want to use them.
// Maybe package are used in every other packages.

// system environment defines
var (
    Name                 string
    Environment          string
    StartTime            time.Time
    SysPath              string                                 // system base path
    RunPath              string                                 // system run path

    LoggerFactory        map[string]seelog.LoggerInterface      // logger instance factory
    Config               ConfigStruct                           // parsed Config structer
)

// Project consts here
var (
    DemoQueue            chan int
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
    SIGNAL_START SIGNAL  = iota + 1
    SIGNAL_STOP
    SIGNAL_RELOAD
)

func (t SIGNAL) String() string {
    switch t {
    case SIGNAL_START:
        return "Start"
    case SIGNAL_STOP:
        return "Stop"
    case SIGNAL_RELOAD:
        return "Reload"
    default:
        return "unknown"
    }
}


// Module running status.
type RState int
const (
    Running RState       = iota + 1
    Stopped
)

func (s RState) String() string {
    switch s {
    case Running:
        return "Running"
    case Stopped:
        return "Stopped"
    default:
        return "Unknown"
    }
}


// Modules monitor log level defines.
const (
    MONITOR_INFO int     = iota + 1
    MONITOR_WARN
    MONITOR_ERROR
    MONITOR_FATAL
)


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
