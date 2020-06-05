package goconf

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

//ReadCsv csv 文件读取
func ReadCsv(path string, message chan []string, f func([]string), fNum int) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	count := 0
	for {
		recoder, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		count++
		if count%fNum == 0 {
			f(recoder)
		}
		message <- recoder
	}
	return nil
}

//ReadCsvWorker 读取文件,并且有worker 消耗文件
//path 文件路径 f 读取函数中间操作, fNum读取多少个文件执行f函数
//worker 工作池 ,workerNumber  工作池数量
func ReadCsvWorker(path string, f func([]string), fNum int, worker func(chan []string, *sync.WaitGroup), workerNumber int) error {
	message := make(chan []string, 10000)
	var wg sync.WaitGroup
	for i := 0; i < workerNumber; i++ {
		wg.Add(1)
		go worker(message, &wg)
	}
	ReadCsv(path, message, f, fNum)
	for i := 0; i < workerNumber; i++ {
		message <- []string{"ok"}
	}
	wg.Wait()
	return nil
}
