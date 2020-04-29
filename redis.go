package goconf

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	redisv7 "github.com/go-redis/redis/v7"
)

//newPool 创建连接池
func newPool(server string, option redis.DialOption, redisMaxActivePool int, redisMaxIdle int, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		MaxActive:   redisMaxActivePool,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, option, option)
			if err != nil {
				msg := fmt.Sprintf("error:%v", err)
				fmt.Println(msg)
				return nil, err
			}
			if _, err := c.Do("select", db); err != nil {
				c.Close()
				msg := fmt.Sprintf("error:%v", err)
				fmt.Println(msg)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

//NewRedisConfig 创建新的redis配置连接池
func newRedisConfig(redisKey string, db int) *redis.Pool {
	sec := nowConfig.Section("redis")
	_redisURL := redisKey + "_redis_url"
	_redisPwd := redisKey + "_redis_pwd"
	_redisMaxIdle := redisKey + "_redis_max_idle"
	_redisMaxActivePool := redisKey + "_redis_max_active_pool"

	redisURL := sec.Key(_redisURL).MustString("")
	redisPwd := sec.Key(_redisPwd).MustString("")
	redisMaxIdle := sec.Key(_redisMaxIdle).MustInt(20)
	redisMaxActivePool := sec.Key(_redisMaxActivePool).MustInt(20)
	option := redis.DialPassword(redisPwd)
	_pool := newPool(redisURL, option, redisMaxActivePool, redisMaxIdle, db)
	return _pool
}

//InitRedis 初始化redis
func InitRedis(key string, db int) *redis.Pool {
	// pool 获取客户端 需要显式close
	pool := newRedisConfig(key, db)
	return pool
}

//InitRedisConfDb 初始化redis
func InitRedisConfDb(redisKey string) *redis.Pool {
	// pool 获取客户端 需要显式close
	sec := nowConfig.Section("redis")
	_redisURL := redisKey + "_redis_url"
	_redisPwd := redisKey + "_redis_pwd"
	_redisMaxIdle := redisKey + "_redis_max_idle"
	_redisMaxActivePool := redisKey + "_redis_max_active_pool"
	_redisDB := redisKey + "_redis_db"

	redisURL := sec.Key(_redisURL).MustString("")
	redisPwd := sec.Key(_redisPwd).MustString("")
	redisMaxIdle := sec.Key(_redisMaxIdle).MustInt(20)
	db := sec.Key(_redisDB).MustInt(-1)
	if db == -1 {
		panic("redis db not config")
	}
	redisMaxActivePool := sec.Key(_redisMaxActivePool).MustInt(20)
	option := redis.DialPassword(redisPwd)
	_pool := newPool(redisURL, option, redisMaxActivePool, redisMaxIdle, db)
	return _pool
}

//InitRedis1ConfDb 初始化redis
func InitRedis1ConfDb(key string) *redisv7.Client {
	sec := nowConfig.Section("redis")
	_redisURL := key + "_redis_url"
	_redisPwd := key + "_redis_pwd"
	_redisDB := key + "_redis_db"
	redisURL := sec.Key(_redisURL).MustString("")
	redisPwd := sec.Key(_redisPwd).MustString("")
	db := sec.Key(_redisDB).MustInt(-1)
	if db == -1 {
		panic("redis db not config")
	}
	client1 := redisv7.NewClient(&redisv7.Options{
		Addr:     redisURL,
		Password: redisPwd,
		DB:       db,
	})
	cmd := client1.Ping()
	if cmd.Err() != nil {
		panic(cmd.Err())
	}
	return client1
}

//InitRedis1 初始化redis
func InitRedis1(key string, db int) *redisv7.Client {
	sec := nowConfig.Section("redis")
	_redisURL := key + "_redis_url"
	_redisPwd := key + "_redis_pwd"
	redisURL := sec.Key(_redisURL).MustString("")
	redisPwd := sec.Key(_redisPwd).MustString("")
	client1 := redisv7.NewClient(&redisv7.Options{
		Addr:     redisURL,
		Password: redisPwd,
		DB:       db,
	})
	cmd := client1.Ping()
	if cmd.Err() != nil {
		panic(cmd.Err())
	}
	return client1
}
