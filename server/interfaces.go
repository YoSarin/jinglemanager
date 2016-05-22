package server

// LogI - intrface for logger
type LogI interface {
	Info(string)
	Warning(string)
	Error(string)
	Debug(string)
	Notice(string)
}

// SongI - interface used to play songs
type SongI interface {
	Play(func())
	Pause()
	Stop()
	IsPlaying() bool
	Position() float64
	FileName() string
}
