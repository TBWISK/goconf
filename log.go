package goconf

import (
	"fmt"
	"io"
	"net"
	"time"

	"log"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//GetLogPath 获取日志路径
func GetLogPath() string {
	sec := nowConfig.Section("log")
	return sec.Key("log_path").MustString("./temp/")
}

//NewLoger 初始化
func NewLoger(path string) *zap.Logger {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(path + "/info.log")
	warnWriter := getWriter(path + "/error.log")

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
	)

	return zap.New(core, zap.AddCaller())
}
func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		// filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

//LoggerItem 容量切割
type LoggerItem struct {
	logger *log.Logger
}

//LoggerTools 日志工具
type LoggerTools struct {
	info *LoggerItem
	err  *LoggerItem
}

//NewLoggerTools log 初始化
func NewLoggerTools(path string) *LoggerTools {
	infoPath := path + "/info.log"
	errPath := path + "/error.log"
	info := NewLoggerItem(infoPath)
	err := NewLoggerItem(errPath)
	return &LoggerTools{info: info, err: err}
}

//INFO info打印
func (c *LoggerTools) INFO(mark string, brief string, msg string) {
	c.info.INFO(mark, brief, msg)
}

//ERROR error打印
func (c *LoggerTools) ERROR(mark string, brief string, msg string) {
	c.err.ERROR(mark, brief, msg)
}

//WARN 告警打印
func (c *LoggerTools) WARN(mark string, brief string, msg string) {
	c.err.WARN(mark, brief, msg)
}

func init() {
	LocationIP = getInternal()
}

//NewLoggerItem 创建日志工具
func NewLoggerItem(Filename string) *LoggerItem {
	logger := log.New(&lumberjack.Logger{
		Filename:   Filename,
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     360, //days
	}, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	return &LoggerItem{logger: logger}

}

//LocationIP 本地ip
var LocationIP string

//Infof 打印普通日志
func (c *LoggerItem) Infof(args ...interface{}) {
	c.logger.Printf("%v %v - %v", LocationIP, "INFO", args)
}

//Error 打印错误日志
func (c *LoggerItem) Error(err interface{}, funcName string, args ...interface{}) {
	item := fmt.Sprintf("%v", args)
	c.logger.Printf("%v %v %v [func_name:%s][args:%s]", "ERROR", LocationIP, err, funcName, item)
}

//Errorf 打印错误日志
func (c *LoggerItem) Errorf(args ...interface{}) {
	c.logger.Printf("%v %v - %s", "ERROR", LocationIP, args)
}

//Warnf 打印错误日志
func (c *LoggerItem) Warnf(args ...interface{}) {
	c.logger.Printf("%v %v - %v", "WARN", LocationIP, args)
}

//INFO info信息打印
func (c *LoggerItem) INFO(mark string, brief string, msg string) {
	line := fmt.Sprintf("INFO %v %v - %v; %v", LocationIP, mark, brief, msg)
	c.logger.Println(line)
}

//WARN warn
func (c *LoggerItem) WARN(mark string, brief string, msg string) {
	line := fmt.Sprintf("WARN %v %v - %v; %v", LocationIP, mark, brief, msg)
	c.logger.Println(time.Now(), line)
}

//ERROR warn
func (c *LoggerItem) ERROR(mark string, brief string, msg string) {
	line := fmt.Sprintf("ERROR %v %v - %v; %v", LocationIP, mark, brief, msg)
	c.logger.Println(time.Now(), line)
}

//getInternal 获取本机ip 多个网卡只获取第一个
func getInternal() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("error no ip")
		return "127.0.0.1"
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// os.Stdout.WriteString(ipnet.IP.String() + "\n")
				ip = ipnet.IP.String()
				return ip
			}
		}
	}
	return ip
}
