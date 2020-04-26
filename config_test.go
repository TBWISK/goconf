package goconf

import (
	"context"
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
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
func Test_ConfigMongoInit(t *testing.T) {
	projectPath := ""
	NewConfigParse(projectPath)
	mgoClient := InitMongo("xxx")
	if err := mgoClient.Ping(context.Background(), nil); err != nil {
		fmt.Println("error", err)
	}
}
