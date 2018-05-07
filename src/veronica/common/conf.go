package common

// This is yaml config file mapping, map to ect/conf.ENVIRONMENT.yaml.
// We used "gopkg.in/yaml.v2" to parse yaml files.
// You should add your own Config into this.
// You can also change config parsers if you wish to. ~_~
//
// @AUTHOR tangyang

type ConfigStruct struct {
    Environ  string
    Worker   int
    ChanSize int            `yaml:"max_chan_size"`
    Log      LogDef
    Debug    DebugDef
    Monitor  MonitorDef
    Demo     DemoDef

    Tickers   map[string]int64
    /*database config*/
    Databases map[string]struct {
        Host     string
        Port     int
        User     string
        Password string
        DbName   string     `yaml:"db_name"`
        LogMode  bool       `yaml:"log_mode"`
    }
    /*redis config*/
    Redis   map[string]struct {
        Host        string
        Db          int
        Password    string
        MaxIdle     int     `yaml:"max_idle"`
        IdleTimeout int     `yaml:"idle_timeout"`
    }
}

/**
 * debug def
 */
type DebugDef struct {
    Debug  int
    Remote string
}

/**
 * monitor def
 */
type MonitorDef struct {
    Interval int
}

/**
 * log def
 */
type LogDef struct {
    Rotate int
    Format int
    Level  int
    Path   string
    IsTerminal int          `yaml:"is_terminal"`
}

type DemoDef struct {
    Threads int
    Test    string
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
