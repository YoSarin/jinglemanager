package lib

// Context - context containing app information
type Context struct {
	Log   LogI
	Songs *FileList
	Sound *SoundController
}

func (c *Context) cleanup() {
	c.Songs = NewFileList()
	c.Sound = NewSoundController(c.Log)
}
