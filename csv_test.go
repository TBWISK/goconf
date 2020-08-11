package goconf

import (
	"fmt"
	"sync"
	"testing"
)

func Test_Csv(t *testing.T) {
	message := make(chan []string, 10000)
	f := func(items []string, cnt int) {
		fmt.Println("xxx", items)
	}
	ReadCsv("/Users/tbwisk/Downloads/imei.csv", message, f, 100)
}
func Test_ReadCsvWorker(t *testing.T) {
	f := func(items []string, cnt int) {
		fmt.Println("xxx", items)
	}
	worker := func(items chan []string, wg *sync.WaitGroup) {
		for {
			obj := <-items
			if len(obj) == 0 {
				break
			}
			if obj[0] == "ok" {
				break
			}
			fmt.Println("obj", obj)
		}
		wg.Done() //wg.Add(1) 默认已经添加
	}
	ReadCsvWorker("/Users/tbwisk/Downloads/imei.csv", f, 1000, worker, 10)
	fmt.Println("finish")
}
