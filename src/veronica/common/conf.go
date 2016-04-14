package common


// This is yaml config file mapping, map to ect/conf.ENVIRONMENT.yaml.
// We used "gopkg.in/yaml.v2" to parse yaml files.
// You should add your own Config into this.
// You can also change config parsers if you wish to. ~_~
//
// @AUTHOR tangyang
import(
)


type ConfigStruct struct {
    Environment      string

    Global           GlobalConf
    Demo             DemoConf

    Tickers          map[string]int64
    Logger           map[string]string
    Databases map[string]struct {
        Host         string
        Port         int
        User         string
        Password     string
        DbName       string     `yaml:"db_name"`
        LogMode      bool       `yaml:"log_mode"`
    }
    Redis map[string]struct {
        Host         string
        Db           int
        Password     string
        MaxIdle      int        `yaml:"max_idle"`
        IdleTimeout  int        `yaml:"idle_timeout"`
    }
}

type GlobalConf struct {
    Processors       int
    MaxChannelSize   int        `yaml:"channel_size"`
    MonitorInterval  int        `yaml:"monitor_interval"`

    PprofMode        bool       `yaml:"pprof_mode"`
    PprofAddr        string     `yaml:"pprof_addr"`
}

type DemoConf struct {
    Threads          int
    Test             string
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
