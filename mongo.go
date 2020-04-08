package goconf

import (
	"fmt"

	"github.com/globalsign/mgo"
)

// 初始化
func newMongoConfig(key string) *mgo.Session {
	sec := nowConfig.Section("mongo")
	mongodbURL := key + "_mongo_url"
	mongoIsAuth := key + "_mongo_is_auth"
	mongoUser := key + "_mongo_user"
	mongoPassword := key + "_mongo_password"
	// 判断是否为空

	mgoSession, err := mgo.Dial(sec.Key(mongodbURL).MustString(""))
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Eventual, true)
	isAuth := sec.Key(mongoIsAuth).MustInt(0)
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

//MongoInit 对mongodb的初始化
func MongoInit(key string) *mgo.Session {
	// mongodb 初始化
	return newMongoConfig(key)
}
