package lib

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

// Context - context containing app information
type Context struct {
	Log        LogI
	Songs      *UniqueList
	Sound      *SoundController
	Jingles    *UniqueList
	Tournament *Tournament
}

// NewTournament - will prepare context for new tournament
func (c *Context) NewTournament(name string) {
	c.cleanup()
	c.Tournament.SetName(name)
}

func (c *Context) cleanup() {
	c.Songs = NewUniqueList()
	c.Sound = NewSoundController(c.Log)
	c.Tournament = NewTournament("")
	c.Jingles = NewUniqueList()
	ChannelChange.Emit(EventTypeCleanup, struct{}{})
}

// StorageDir - return path to current tournament directory (and creates path if necessarry)
func (c *Context) StorageDir() string {
	u, _ := user.Current()
	p := path.Join(u.HomeDir, ".jinglemanager", c.Tournament.Name, "media")
	os.MkdirAll(p, os.ModeDir)
	return path.Join(u.HomeDir, ".jinglemanager", c.Tournament.Name)
}

// AppDir - return path to application directory
func (c *Context) AppDir() string {
	u, _ := user.Current()
	p := path.Join(u.HomeDir, ".jinglemanager")
	os.MkdirAll(p, os.ModeDir)
	return path.Join(u.HomeDir, ".jinglemanager")
}

// LastTournament - return path to application directory
func (c *Context) LastTournament() string {
	f, _ := os.Open(path.Join(c.AppDir(), "last.tournament"))
	t, _ := ioutil.ReadAll(f)
	return string(t)
}
