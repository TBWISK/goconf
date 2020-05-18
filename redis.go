package goconf

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	redisv7 "github.com/go-redis/redis/v7"
)

//RedisConf redis配置信息
type RedisConf struct {
	URL           string
	Pwd           string
	MaxIdle       int
	MaxActivePool int
	DB            int
}

//NewRedisConf redis配置初始化
func NewRedisConf(key string) *RedisConf {
	section := nowConfig.Section("redis")
	_redisURL := key + "_redis_url"
	_redisPwd := key + "_redis_pwd"
	_redisMaxIdle := key + "_redis_max_idle"
	_redisMaxActivePool := key + "_redis_max_active_pool"
	_redisDB := key + "_redis_db"
	redisURL := section.Key(_redisURL).MustString("")
	redisPwd := section.Key(_redisPwd).MustString("")
	redisMaxIdle := section.Key(_redisMaxIdle).MustInt(5)
	redisMaxActivePool := section.Key(_redisMaxActivePool).MustInt(20)
	db := section.Key(_redisDB).MustInt(-1)
	return &RedisConf{
		URL:           redisURL,
		Pwd:           redisPwd,
		MaxIdle:       redisMaxIdle,
		MaxActivePool: redisMaxActivePool,
		DB:            db,
	}
}

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

//InitRedis 初始化redis
func InitRedis(key string, db int) *redis.Pool {
	// pool 获取客户端 需要显式close
	conf := NewRedisConf(key)
	option := redis.DialPassword(conf.Pwd)
	pool := newPool(conf.URL, option, conf.MaxActivePool, conf.MaxIdle, db)
	return pool
}

//InitRedisConfDb 初始化redis
func InitRedisConfDb(key string) *redis.Pool {
	// pool 获取客户端 需要显式close
	conf := NewRedisConf(key)
	option := redis.DialPassword(conf.Pwd)
	db := conf.DB
	if db == -1 {
		panic("redis db not config")
	}
	pool := newPool(conf.URL, option, conf.MaxActivePool, conf.MaxIdle, db)
	return pool
}

//InitRedis1ConfDb 初始化redis
func InitRedis1ConfDb(key string) *redisv7.Client {
	conf := NewRedisConf(key)
	if conf.DB == -1 {
		panic("redis db not config")
	}
	client1 := redisv7.NewClient(&redisv7.Options{
		Addr:         conf.URL,
		Password:     conf.Pwd,
		DB:           conf.DB,
		PoolSize:     conf.MaxActivePool,
		MinIdleConns: conf.MaxIdle,
	})
	cmd := client1.Ping()
	if cmd.Err() != nil {
		panic(cmd.Err())
	}
	return client1
}

//InitRedis1 初始化redis
func InitRedis1(key string, db int) *redisv7.Client {
	conf := NewRedisConf(key)
	client1 := redisv7.NewClient(&redisv7.Options{
		Addr:         conf.URL,
		Password:     conf.Pwd,
		DB:           db,
		PoolSize:     conf.MaxActivePool,
		MinIdleConns: conf.MaxIdle,
	})
	cmd := client1.Ping()
	if cmd.Err() != nil {
		panic(cmd.Err())
	}
	return client1
}
