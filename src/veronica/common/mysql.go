package common


// Get MySQL connections.
// Pkg use "github.com/jinzhu/gorm" & "github.com/go-sql-driver/mysql".
// Struct parse dbFlag and get belong connections.
// @author tangyang

import (
    "fmt"

    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
)

type DbClient struct {
    Conn *gorm.DB
}

func NewDbClient(name string) (*DbClient, error) {
    cluster, ok := Config.DbClusters[name]
    if !ok {
        return nil, fmt.Errorf("db cluster: %s, not exist", name)
    }

    host    := cluster.Host
    port    := cluster.Port
    user    := cluster.User
    passwd  := cluster.Password
    dbName  := cluster.DbName
    logMode := cluster.LogMode
    fd      := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd, host, port, dbName)
    conn, err := gorm.Open("mysql", fd)
    if err != nil {
        return nil, err
    }
    // log detail log or false
    if logMode {
        conn.LogMode(true)
    } else {
        conn.LogMode(false)
    }
    return &DbClient{
        Conn : conn,
    }, nil
}

/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
