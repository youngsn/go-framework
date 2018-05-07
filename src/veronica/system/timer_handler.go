package system


// Trigger task handler defines.
// @AUTHOR tangyang
import (
    "fmt"
    "math/rand"

    c "veronica/common"
)

type TickerHandler func()

// Get task handler func handle.
func getTickerHandler(tickerName string) TickerHandler {
    var handler TickerHandler
    switch {
    case tickerName == "demoDoing":             // THIS IS JUST DEMO
        handler = demoDoing
    default:
        panic(fmt.Sprint("No right handler to TimerTask: ", tickerName))
    }
    return handler
}

// THIS IS JUST demo handler func.
func demoDoing() {
    num := rand.Intn(100)
    c.Logger.Infof("Send %d to DemoQueue", num)
    c.DemoQueue<- num
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
