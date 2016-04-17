package lib

// Context - context containing app information
type Context struct {
	Log        LogI
	Songs      *FileList
	Sound      *SoundController
	Tournament *Tournament
}

func (c *Context) cleanup() {
	c.Songs = NewFileList()
	c.Sound = NewSoundController(c.Log)
	c.Tournament = NewTournament("")
}
