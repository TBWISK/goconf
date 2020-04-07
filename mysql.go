package goconf

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//newMysqlConfig 配置初始化
func newMysqlConfig(key string) *sql.DB {
	sec := nowConfig.Section("mysql")
	dbURL := sec.Key(key + "_mysql_url").MustString("")
	dbUser := sec.Key(key + "_mysql_user").MustString("")
	dbPw := sec.Key(key + "_mysql_pwd").MustString("")
	dbName := sec.Key(key + "_mysql_database").MustString("")
	dbSubURL := dbUser + ":" + dbPw + "@tcp(" + dbURL + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"
	sqlDb, err := sql.Open("mysql", dbSubURL)
	if err != nil {
		msg := fmt.Sprintf("error:%v", err)
		fmt.Println(msg)
		panic(err)
	}
	return sqlDb
}

//InitMysql 初始化mysql
func InitMysql(key string) *sql.DB {
	db := newMysqlConfig(key)
	return db
}
