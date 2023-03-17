package goconf

import (
	"fmt"
	"testing"
)

func Test_log(t *testing.T) {
	projectPath := "/Users/tbwisk/coding/github/goconf"
	NewConfigParse(projectPath)
	path := GetLogPath()
	fmt.Println(path)
	logger := NewLoger(path)
	sugar := logger.Sugar()
	sugar.Info("xxx", "xxx")
	sugar.Errorw("msg", "keysAndValues", "2")
	sugar.Warn("debug")
	item := map[string]interface{}{"xxx": 1, "hjelo": "world"}
	sugar.Info(item)

}
