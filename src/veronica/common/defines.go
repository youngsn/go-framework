package common


// This file used to add public defines.
// Like public struct defines, public interfaces, etc.
// @author tangyang


// Public interface defines.

// System module interface define.
// If you want to use public libs, you Module need implements these methods.
type Module interface {
    Status()  RState                  // module running status
    Monitor() []*MonitorPack          // stdout module status to monitor
    Ctrl(s <-chan SIGNAL)             // main thread to send signal into module chan
}


// Monitor data structure.
// StdLevel used to define monitor logging level.
// Content is logging content.
type MonitorPack struct {
    StdLevel      int       // monitor level
    Content       string    // monitor content
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
