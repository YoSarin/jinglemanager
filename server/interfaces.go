package server

// LogI - intrface for logger
type LogI interface {
	Info(string)
	Warning(string)
	Error(string)
	Debug(string)
	Notice(string)
}
