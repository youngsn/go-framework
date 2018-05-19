package common


// This file used to add public defines.
// Like public struct defines, public interfaces, etc.
// @author tangyang
import (
    "fmt"
)

// Module type define
// If you want to use veronica, your own module should MUST implement these methods.
type Module interface {
    Init(s <-chan SIGNAL)          // init module, listen signal chan
    Status()  RState               // module exec status
    Monitor() []*MonitorPack       // monitor module status
}

// Monitor pack data.
type MonitorPack struct {
    State   RState                  // module run status
    Level   int                     // log level
    Content string                  // log content
    Fields  LogFields               // log fields
}

type NameServiceInst struct {
    Server       *Server
    Retry        int
    ConnTimeout  int64
    ReadTimeout  int64
    WriteTimeout int64
}

type Server struct {
    Host string
    Port int
}

func (s Server) String() string {
    return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
