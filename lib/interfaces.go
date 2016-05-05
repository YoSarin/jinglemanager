package lib

// LogI - intrface for logger
type LogI interface {
	Info(string, ...interface{})
	Warning(string, ...interface{})
	Error(string, ...interface{})
	Debug(string, ...interface{})
	Notice(string, ...interface{})
}
