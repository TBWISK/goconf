package goconf

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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

//GetInsertSQL 获取sql
func GetInsertSQL(table string, sets map[string]interface{}) string {
	keys := []string{}
	values := []string{}
	for idx := range sets {
		value := sets[idx]
		switch value.(type) {
		case string:
			keys = append(keys, idx)
			values = append(values, "'"+value.(string)+"'")
		case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, float32, float64:
			keys = append(keys, idx)
			var x string
			switch value.(type) {
			case int:
				x = strconv.FormatInt(int64(value.(int)), 10)
			case int16:
				x = strconv.FormatInt(int64(value.(int16)), 10)
			case int32:
				x = strconv.FormatInt(int64(value.(int32)), 10)
			case int64:
				x = strconv.FormatInt(int64(value.(int64)), 10)
			case int8:
				x = strconv.FormatInt(int64(value.(int8)), 10)
			case uint:
				x = strconv.FormatUint(uint64(value.(uint)), 10)
			case uint16:
				x = strconv.FormatUint(uint64(value.(uint16)), 10)
			case uint32:
				x = strconv.FormatUint(uint64(value.(uint32)), 10)
			case uint64:
				x = strconv.FormatUint(uint64(value.(uint64)), 10)
			case float32:
				x = strconv.FormatFloat(float64(value.(float32)), 'f', -6, 32)
			case float64:
				x = strconv.FormatFloat(float64(value.(float64)), 'f', -6, 64)
			}
			values = append(values, x)
		default:
			panic("GetInsertSQL type error")
		}
	}
	fmt.Println(keys)
	fmt.Println(values)
	sql := fmt.Sprintf("insert into %v (%v) values (%v)", table, strings.Join(keys, ","), strings.Join(values, ","))
	return sql
}

//ExecuteSQL sql转换
func ExecuteSQL(table string, sets map[string]interface{}, db *sql.DB) (int64, error) {
	keys := []string{}
	values := []string{}
	objs := []interface{}{}
	for idx := range sets {
		value := sets[idx]
		switch value.(type) {
		case string:
			keys = append(keys, idx)
			values = append(values, "?")
		case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, float32, float64:
			keys = append(keys, idx)
			values = append(values, "?")
		default:
			panic("GetInsertSQL type error")
		}
		objs = append(objs, value)
	}
	sqlx := fmt.Sprintf("insert into %v (%v) values (%v)", table, strings.Join(keys, ","), strings.Join(values, ","))
	return Execute(table, sqlx, objs, db)
}

//Execute sql 执行
func Execute(table string, sql string, values []interface{}, db *sql.DB) (int64, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err, "Prepare")
		return 0, err
	}
	result, err := stmt.Exec(values...)
	if err != nil {
		fmt.Println(err, "Exec")
		return 0, err
	}
	return result.LastInsertId()
}

//FormatUpdateSQL 整理uppdate的sql语句
func FormatUpdateSQL(sets map[string]interface{}, where map[string]interface{}) {

}
