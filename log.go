package gonet

type ILogger interface {
	Info(args ...any)
	Warn(args ...any)
	Debug(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
