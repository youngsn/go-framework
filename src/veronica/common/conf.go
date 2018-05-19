package common

// This is yaml config file mapping, map to ect/conf.ENVIRONMENT.yaml.
// We used "gopkg.in/yaml.v2" to parse yaml files.
// You should add your own Config into this.
// You can also change config parsers if you wish to. ~_~
//
// @AUTHOR tangyang

type AppConfig struct {
    Environ  string
    Worker   int
    ChanSize int            `yaml:"max_chan_size"`
    Debug    struct {Debug int; Addr string}
    Monitor  struct {Interval int}
    Log      struct {
        Rotate int
        Level  int
        Path   string
    }
    Tickers  map[string]int64
    Service  struct {
        Local NameLocalDef
    }
    Redis    map[string]struct {
        Db      int
        Passwd  string
        Service string
    }
    DbClusters map[string]DbClusterDef

    Demo DemoDef
}

type DbClusterDef struct {
    Host     string
    Port     int
    User     string
    Password string
    DbName   string     `yaml:"db_name"`
    LogMode  bool       `yaml:"log_mode"`
}

type DemoDef struct {
    Threads int
    Test    string
}

type NameLocalDef map[string]struct {
    Port         int            `yaml:"port"`
    Retry        int            `yaml:"retry"`
    ConnTimeout  int64          `yaml:"ctimeout"`
    ReadTimeout  int64          `yaml:"rtimeout"`
    WriteTimeout int64          `yaml:"wtimeout"`
    Strategy     struct {
        Balance  string
    }
    Servers      []Server
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
