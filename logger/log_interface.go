package logger

// 公用接口，所调用对象皆可使用实例
type LogInterface interface {
	Init()
	SetLevel(level int)
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Close()
}
