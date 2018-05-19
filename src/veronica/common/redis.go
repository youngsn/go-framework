package common

import(
    "fmt"
    "time"

    "github.com/go-redis/redis"
)

type Redis struct {
    Client *redis.Client
}

func NewRedis(cluster string) (*Redis, error) {
    inst     := &Redis{}
    cli, err := inst.Connect(cluster)
    if err != nil {
       return nil, err
    }
    inst.Client = cli
    return inst, err
}

func (r *Redis) Connect(cluster string) (*redis.Client, error) {
    if _, ok := Config.Redis[cluster]; !ok {
        return nil, fmt.Errorf("redis cluster: %s, not exist", cluster)
    }

    conf     := Config.Redis[cluster]                       // get cluster config
    ret, err := NameCli.GetService(conf.Service)            // get service config from name client
    if err != nil {
        return nil, err
    }
    db       := conf.Db
    passwd   := conf.Passwd
    addr     := ret.Server.String()
    retry    := ret.Retry
    cTimeout := time.Duration(ret.ConnTimeout) * time.Millisecond
    rTimeout := time.Duration(ret.ReadTimeout) * time.Millisecond
    wTimeout := time.Duration(ret.WriteTimeout) * time.Millisecond

    client   := redis.NewClient(&redis.Options{
        Addr         : addr,
        DB           : db,
        Password     : passwd,
        MaxRetries   : retry,
        DialTimeout  : cTimeout,
        ReadTimeout  : rTimeout,
        WriteTimeout : wTimeout,
    })
    if _, err := client.Ping().Result(); err != nil {       // test connection
        return nil, err
    }
    Logger.WithFields(LogFields{
        "cluster"  : cluster,
        "addr"     : addr,
        "retry"    : retry,
        "ctimeout" : cTimeout,
        "rtimeout" : rTimeout,
        "wtimeout" : wTimeout,
    }).Infof("connect redis: %s, success", cluster)
    return client, nil
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
