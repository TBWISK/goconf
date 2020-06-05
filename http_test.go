package goconf

import (
	"fmt"
	"testing"
)

func Test_Get(t *testing.T) {
	url := "http://baidu.com"
	resp, err := Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resp))
}
