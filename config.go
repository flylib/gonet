package goNet

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// 设置日志格式为json格式
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,                  //显示颜色
		FullTimestamp:   true,                  //显示时间
		TimestampFormat: "2006/01/02 15:04:05", //配置时间显示格式
	})
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	logrus.SetOutput(os.Stdout)
	// 设置日志级别为warn以上
	logrus.SetLevel(logrus.InfoLevel)
}
