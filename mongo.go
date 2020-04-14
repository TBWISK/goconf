package goconf

import (
	"context"
	"fmt"
	"strings"

	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 初始化
func newMongoConfig(key string) *mgo.Session {
	sec := nowConfig.Section("mongo")
	mongodbURL := key + "_mongo_url"
	mongoIsAuth := key + "_mongo_auth"
	mongoUser := key + "_mongo_user"
	mongoPassword := key + "_mongo_password"
	// 判断是否为空
	mgoSession, err := mgo.Dial(sec.Key(mongodbURL).MustString(""))
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Eventual, true)
	isAuth := sec.Key(mongoIsAuth).MustInt(0)
	fmt.Println("isAuth", isAuth)
	if isAuth == 1 {
		root := sec.Key(mongoUser).MustString("")
		Password := sec.Key(mongoPassword).MustString("")
		credential := &mgo.Credential{Username: root, Password: Password}
		err = mgoSession.Login(credential)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	return mgoSession
}

//MgoInit 对mongodb的初始化 第三方库
func MgoInit(key string) *mgo.Session {
	// mongodb 初始化
	return newMongoConfig(key)
}

//MongoInit 初始化 官方库;暂时只支持单个mongo
func MongoInit(key string) *mongo.Client {
	sec := nowConfig.Section("mongo")
	mongodbURL := key + "_mongo_url"
	mongoUser := key + "_mongo_user"
	mongoPassword := key + "_mongo_password"
	mongoIsAuth := key + "_mongo_is_auth"
	opts := options.Client()
	hosts := strings.Split(sec.Key(mongodbURL).MustString(""), ",")
	Username := sec.Key(mongoUser).MustString("")
	Password := sec.Key(mongoPassword).MustString("")
	opts = opts.SetHosts(hosts)
	if sec.Key(mongoIsAuth).MustInt(0) == 1 {
		opts.SetAuth(options.Credential{
			Username: Username,
			Password: Password})
	}
	client, err := mongo.NewClient(opts)
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	return client
}
