package system


// Trigger task handler defines.
// @AUTHOR tangyang

import (
    "fmt"
    "math/rand"

    . "veronica/common"

    Log "github.com/cihub/seelog"
)


type TickerHandler func()

// Get task handler func handle.
func getTickerHandler(tickerName string) TickerHandler {
    var handler TickerHandler
    switch {
    case tickerName == "demoDoing":             // THIS IS JUST DEMO
        handler      = demoDoing
    default:
        panic(fmt.Sprint("No right handler to TimerTask: ", tickerName))
    }
    return handler
}

// THIS IS JUST demo handler func.
func demoDoing() {
    num := rand.Intn(100)
    Log.Infof("Send %d to DemoQueue", num)
    DemoQueue<- num
    Log.Infof("Send success")
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
