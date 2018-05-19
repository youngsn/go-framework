package common

import(
    "fmt"
    "time"
    "math/rand"
)

var (
    DEF_PORT  int = 8099
    DEF_RETRY int = 0
    DEF_CTIMEOUT int64 = 500
    DEF_RTIMEOUT int64 = 500
    DEF_WTIMEOUT int64 = 500

    NAME_CACHE_TIME int64 = 60         // cache time, second
)

type NameClient struct {
    nameCache map[string]*nameCache
}

// name service cache
type nameCache struct {
    Ts      time.Time
    Inst    *serviceInst
    Servers []*Server
}

func NewNameClient() *NameClient {
    return &NameClient{
        nameCache : map[string]*nameCache{},
    }
}

func (n *NameClient) GetService(clusterName string) (*NameServiceInst, error) {
    localConf := Config.Service.Local

    var (
        err error = fmt.Errorf("cluster: %s, config not exist", clusterName)
        srvInst *serviceInst
        servers []*Server
    )
    if _, ok := localConf[clusterName]; ok {
        name := clusterName
        srvInst, servers, err = n.getLocalInst(name)
    }
    if err != nil {
        return nil, err
    }

    // 根据策略选择服务地址
    var host *Server
    if srvInst.Strategy.Balance == "Random" {
        host = n.randomServers(servers)
    } else {
        host = n.randomServers(servers)
    }
    return &NameServiceInst{
        Server       : host,
        Retry        : srvInst.Retry,
        ConnTimeout  : srvInst.ConnTimeout,
        ReadTimeout  : srvInst.ReadTimeout,
        WriteTimeout : srvInst.WriteTimeout,
    }, nil
}

/**
 * 随机挑选访问地址策略
 *
 * @param []*Server     servers
 * @return *Server
 */
func (n *NameClient) randomServers(servers []*Server) *Server {
    slen   := len(servers)
    if slen == 1 {
        return servers[0]
    }
    rand.Seed(time.Now().UnixNano())
    ret    := rand.Intn(slen - 1)           // 随机rand下标，0 ~ len - 1
    return servers[ret]
}

/**
 * 解析本地配置，获取服务和机器实例配置
 */
func (n *NameClient) getLocalInst(name string) (*serviceInst, []*Server, error) {
    locConf, _ := Config.Service.Local[name]
    if len(locConf.Servers) == 0 {
        return nil, nil, fmt.Errorf("local: %s, no valid servers", name)
    }

    var (
        locInst *serviceInst
        servers []*Server
    )
    // 获取本地机器实例列表
    defaultPort := locConf.Port             // 默认端口，未配置端口将使用默认端口
    for _, server := range locConf.Servers {
        if server.Port == 0 {
            server.Port = defaultPort
        }
        servers = append(servers, &server)
    }
    // 获取本地服务配置
    locInst = &serviceInst{
        Port         : DEF_PORT,
        Retry        : DEF_RETRY,
        ConnTimeout  : DEF_CTIMEOUT,
        ReadTimeout  : DEF_RTIMEOUT,
        WriteTimeout : DEF_WTIMEOUT,
        Strategy     : struct{Balance string}{"Random"},
    }
    if locConf.Port > 0 {
        locInst.Port  = locConf.Port
    }
    if locConf.Retry >= -1 {
        locInst.Retry = locConf.Retry
    }
    if locConf.ConnTimeout >= 0 {
        locInst.ConnTimeout = locConf.ConnTimeout
    }
    if locConf.ReadTimeout >= -1 {
        locInst.ReadTimeout = locConf.ReadTimeout
    }
    if locConf.WriteTimeout >= -1 {
        locInst.WriteTimeout = locConf.WriteTimeout
    }
    if locConf.Strategy.Balance == "Random" {
        locInst.Strategy.Balance = locConf.Strategy.Balance
    }
    return locInst, servers, nil
}

// 服务配置实例数据结构
type serviceInst struct {
    Port         int
    Retry        int
    ConnTimeout  int64
    ReadTimeout  int64
    WriteTimeout int64
    Strategy     struct {
        Balance  string
    }
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
