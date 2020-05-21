package goconf

import (
	"context"
	"fmt"
	"testing"
)

func Test_ConfigMongoInit(t *testing.T) {
	projectPath := "/Users/tbwisk/coding/github/goconf"
	NewConfigParse(projectPath)
	mgoClient := InitMongo("xxx")
	if err := mgoClient.Ping(context.Background(), nil); err != nil {
		fmt.Println("error", err)
	}
	c := mgoClient.Database("dmp").Collection("wax")
	query := c.FindOne(context.Background(), map[string]string{"_id": "0000acfbb7a666c807c1c590759008bf"})
	var result map[string]string
	err := query.Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
func Test_ConfigMgoInit(t *testing.T) {
	projectPath := "/Users/tbwisk/coding/github/goconf"
	NewConfigParse(projectPath)
	mgoClient := InitMgo("xxx")
	err := mgoClient.Ping()
	if err != nil {
		panic(err)
	}
	c := mgoClient.DB("dmp").C("wax")
	query := c.FindId("0000acfbb7a666c807c1c590759008bf")
	var result []map[string]string
	err = query.All(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
func TestMgoInsert(t *testing.T) {
	projectPath := "/Users/tbwisk/coding/github/goconf"
	NewConfigParse(projectPath)
	mgoClient := InitMgo("xxx")
	err := mgoClient.Ping()
	if err != nil {
		panic(err)
	}
	c := mgoClient.DB("dmp").C("wax")
	err = c.Insert(map[string]string{"_id": "temp", "pic": "this a pic"})
	if err != nil {
		panic(err)
	}
	err = c.RemoveId("temp")
	if err != nil {
		panic(err)
	}
}
