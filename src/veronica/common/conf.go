package common


// This is toml config file mapping, map to ect/conf.toml.
// We used "github.com/BurntSushi/toml" to parse toml files.
// You should add your own Config into this.
// You can also change config parsers if you wish to. ~_~
// @author tangyang
import(
)


type ConfigStruct struct {
    Name             string
    Environment      string

    Global           GlobalConf
    Demo             DemoConf

    Tickers          map[string]TickerConf
    Databases        map[string]DatabaseConf
    Redis            map[string]RedisConf
    Logger           map[string]LoggerConf
}

type GlobalConf struct {
    MaxProcs         int
    MaxChannelSize   int64
    MonitorInterval  int64

    PprofMode        bool
    PprofAddr        string
}

type DemoConf struct {
    Threads          int
    Test             string
}

type TickerConf struct {
    Interval         int64
}

type DatabaseConf struct {
    Host             string
    Port             int
    User             string
    Password         string
    DbName           string
    LogMode          bool
}

type RedisConf struct {
    Host             string
    DbName           int
    Password         string
}

type LoggerConf struct {
    Conf             string
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
