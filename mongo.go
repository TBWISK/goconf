package goconf

import (
	"context"
	"fmt"
	"strings"

	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoConf mongo 配置文件
type MongoConf struct {
	MongoURL string
	Auth     int
	User     string
	Pwd      string
}

//NewMongoConf 获取mongo conf
func NewMongoConf(key string) *MongoConf {
	sec := nowConfig.Section("mongo")
	mongodbURL := key + "_mongo_url"
	mongoIsAuth := key + "_mongo_auth"
	mongoUser := key + "_mongo_user"
	mongoPassword := key + "_mongo_password"

	return &MongoConf{
		MongoURL: sec.Key(mongodbURL).MustString(""),
		Auth:     sec.Key(mongoIsAuth).MustInt(0),
		User:     sec.Key(mongoUser).MustString(""),
		Pwd:      sec.Key(mongoPassword).MustString(""),
	}
}

// 初始化
func newMongoConfig(key string) *mgo.Session {
	conf := NewMongoConf(key)
	// 判断是否为空
	mgoSession, err := mgo.Dial(conf.MongoURL)
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Eventual, true)
	if conf.Auth == 1 {
		credential := &mgo.Credential{Username: conf.User, Password: conf.Pwd}
		err = mgoSession.Login(credential)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	return mgoSession
}

//InitMgo 对mongodb的初始化 第三方库
func InitMgo(key string) *mgo.Session {
	// mongodb 初始化
	return newMongoConfig(key)
}

//InitMongo 初始化 官方库;暂时只支持单个mongo
func InitMongo(key string) *mongo.Client {
	conf := NewMongoConf(key)
	opts := options.Client()
	hosts := strings.Split(conf.MongoURL, ",")
	opts = opts.SetHosts(hosts)
	if conf.Auth == 1 {
		opts.SetAuth(options.Credential{
			Username: conf.User,
			Password: conf.Pwd})
	}
	client, err := mongo.NewClient(opts)
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	return client
}
