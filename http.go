package goconf

import (
	"io/ioutil"
	"net/http"
)

//Get 用途，封装好对应的get 和post请求
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

//Post 提交
func Post(url string) ([]byte, error) {
	// http.PostForm(url, data)
	// http.Post(url, contentType, body)
	return nil, nil
}
