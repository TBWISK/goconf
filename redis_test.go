package goconf

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func Test_Redis(t *testing.T) {
	projectPath := "/Users/tbwisk/coding/github/goconf"
	NewConfigParse(projectPath)
	pool := InitRedis("xxx", 1)
	con := pool.Get()
	defer con.Close()
	con.Do("set", "key", "value test")
	fmt.Println(redis.String(con.Do("get", "key")))

}

func Test_RedisMock(t *testing.T) {
	// 测试使用的时候，需要注意一下 命令的大小写需要一致,内部实现可能使用了map来处理
	c := redigomock.NewConn()
	c.Command("HGETALL", "person:1").ExpectMap(map[string]string{
		"name": "hello",
		"age":  "42",
	})
	x := GetPerson(c, "1")
	fmt.Println("x=", x)
}
func GetPerson(conn redis.Conn, id string) map[string]string {
	values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("person:%s", id)))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	resp, err := redis.StringMap(values, err)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return resp
}
