package goconf

import (
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
