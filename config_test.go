package goconf

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

func Test_Config(t *testing.T) {
	projectPath := ""
	NewConfigParse(projectPath)
	pool := InitRedis("xxx", 1)
	con := pool.Get()
	defer con.Close()
	con.Do("set", "key", "value test")
	fmt.Println(redis.String(con.Do("get", "key")))
}

type User struct {
	gorm.Model
	Name string
	Age  int64
}

func Test_GORM(t *testing.T) {
	projectPath := "/Users/tbwisk/coding/github/goconf"
	NewConfigParse(projectPath)
	db := InitGorm("xxx")
	db.AutoMigrate(User{})
	user := User{Name: "xxx", Age: 65}
	x := db.Create(&user)
	if x.Error != nil {
		fmt.Println(x.Error)
	}
}
