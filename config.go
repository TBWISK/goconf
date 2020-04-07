package goconf

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

//ConfigParse 配置环境
type ConfigParse struct {
	Path      string
	NowConfig *ini.File
	Env       string
}

//NewConfigParse 初始化
func NewConfigParse(Path string) *ConfigParse {
	if Path[len(Path)-1:len(Path)] == "/" {
		Path = Path[0 : len(Path)-1]
	}
	return &ConfigParse{Path: Path}
}

func (c *ConfigParse) getFilePath(name string) (path string) {
	// pwd, _ := os.Getwd()
	if name == "" {
		path = c.Path + "/resource/app.conf"
	} else {
		path = c.Path + fmt.Sprintf("/resource/app-%v.conf", name)
	}
	return path
}

func (c *ConfigParse) getFileName() string {
	base := c.getFilePath("")
	selfcfg, err := ini.Load(base)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// c.selfcfgFile = selfcfg
	dev := selfcfg.Section("app").Key("app").MustString("")
	c.Env = dev
	if dev == "" {
		fmt.Println("app.conf的app为空")
		os.Exit(1)
	}
	return dev
}

//Init 内部初始化
func (c *ConfigParse) Init() {
	// name:=
	path := c.getFilePath(c.getFileName())
	cf, err := ini.Load(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.NowConfig = cf
}

//GetConfig 获取配置文件
func (c *ConfigParse) GetConfig() *ini.File {
	return c.NowConfig
}

// how to use
// parse := NewConfigParse("project_path")
// parse.Init()
// iniconf = parse.GetConfig()
