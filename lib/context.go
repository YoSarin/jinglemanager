package lib

// Context - context containing app information
type Context struct {
	Log     LogI
	Songs   *FileList
	Sound   *SoundController
	Changes chan interface{}
}

func (c *Context) cleanup() {
	c.Songs = NewFileList()
	c.Sound = NewSoundController(c.Log)
}
