package goNet

import (
	"github.com/Quantumoffices/beego/logs"
)

func init() {
	//日志输出
	logs.SetLogger(logs.AdapterConsole)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
}
