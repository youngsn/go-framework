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


type DbConnection struct {
    Db          *gorm.DB
}


func NewDbConnection(dbFlag string) (*DbConnection, error) {
    if _, exist := Config.Databases[dbFlag]; !exist {
       return nil, fmt.Errorf("flag %s, not exist", dbFlag)
    }

    dbConfig    := Config.Databases[dbFlag]

    host        := dbConfig.Host
    port        := dbConfig.Port
    user        := dbConfig.User
    password    := dbConfig.Password
    dbName      := dbConfig.DbName
    logMode     := dbConfig.LogMode

    dbSource    := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, host, port, dbName)
    dbConn, err := gorm.Open("mysql", dbSource)
    if err != nil {
        return nil, err
    }

    if true == logMode {
        dbConn.LogMode(true)
    } else {
        dbConn.LogMode(false)
    }

    return &DbConnection{
        Db: &dbConn,
    }, nil
}


/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
