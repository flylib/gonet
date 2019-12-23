package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var Log = logrus.New()

func init() {
	Log.Formatter = &logrus.TextFormatter{
		ForceColors:     true,                  //显示颜色
		FullTimestamp:   true,                  //显示时间
		TimestampFormat: "2006/01/02 15:04:05", //配置时间显示格式
	}
	//打印文件名行号
	Log.AddHook(NewLogHook())
	//输出到终端
	Log.Out = os.Stdout
}

// ContextHook for log the call context
type logHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

func NewLogHook(levels ...logrus.Level) logrus.Hook {
	hook := logHook{
		Field:  "[goNet]",
		Skip:   5,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels implement levels
func (hook logHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (hook logHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
