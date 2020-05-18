package goconf

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//MysqlConf 配置
type MysqlConf struct {
	DbURL    string
	LifeTime int
	IdleConn int
	OpenConn int
}

//NewMysqlConf mysql配置
func NewMysqlConf(key string) *MysqlConf {
	sec := nowConfig.Section("mysql")
	dbURL := sec.Key(key + "_mysql_url").MustString("")
	dbUser := sec.Key(key + "_mysql_user").MustString("")
	dbPw := sec.Key(key + "_mysql_pwd").MustString("")
	dbName := sec.Key(key + "_mysql_database").MustString("")
	dbSubURL := dbUser + ":" + dbPw + "@tcp(" + dbURL + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"
	LifeTime := sec.Key(key + "_mysql_conn_life_time").MustInt(300)
	IdleConn := sec.Key(key + "_mysql_idle_conn").MustInt(1)
	OpenConn := sec.Key(key + "_mysql_open_conn").MustInt(10)
	return &MysqlConf{
		DbURL:    dbSubURL,
		LifeTime: LifeTime,
		IdleConn: IdleConn, OpenConn: OpenConn,
	}
}

//InitMysql 初始化mysql
func InitMysql(key string) *sql.DB {
	conf := NewMysqlConf(key)
	dbSubURL := conf.DbURL
	sqlDb, err := sql.Open("mysql", dbSubURL)
	if err != nil {
		msg := fmt.Sprintf("error:%v", err)
		fmt.Println(msg)
		panic(err)
	}
	sqlDb.SetMaxIdleConns(conf.IdleConn)
	sqlDb.SetMaxOpenConns(conf.OpenConn)
	sqlDb.SetConnMaxLifetime(time.Second * time.Duration(conf.LifeTime))
	return sqlDb
}

//InitGorm 初始化
func InitGorm(key string) *gorm.DB {
	conf := NewMysqlConf(key)
	dbSubURL := conf.DbURL
	db, err := gorm.Open("mysql", dbSubURL)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(conf.IdleConn)
	db.DB().SetMaxOpenConns(conf.OpenConn)
	db.DB().SetConnMaxLifetime(time.Second * time.Duration(conf.LifeTime))
	return db
}
