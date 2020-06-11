package goconf

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

//PostJSON 提交post json
func PostJSON(url string, b []byte) ([]byte, error) {
	contentType := "application/json;charset=utf-8"
	body := bytes.NewBuffer(b)
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//PostData 提交post json,b= name=cjb
//url(url,"name=cjb")
func PostData(url string, b string) ([]byte, error) {
	contentType := "application/x-www-form-urlencoded"
	body := strings.NewReader(b)
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//PostForm 提交post 表单
func PostForm(urlx string, b url.Values) ([]byte, error) {
	// body := strings.NewReader(b)
	resp, err := http.PostForm(urlx, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}
