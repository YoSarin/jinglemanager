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

const (
	// EventTypeCleanup - event type related to total cleanup
	EventTypeCleanup = EventType("cleanup")
)

var (
	context        *Context
	ChannelCleanup = Channel{name: "cleanup", allowed: map[EventType]bool{EventTypeCleanup: true}}
)

func NewContext(log LogI) *Context {
	Ctx := &Context{}

	Ctx.Log = log
	Ctx.Songs = NewUniqueList()
	Ctx.Sound = NewSoundController(log)
	Ctx.Tournament = NewTournament("", Ctx)
	Ctx.Jingles = NewUniqueList()

	return Ctx
}

// NewTournament - will prepare context for new tournament
func (c *Context) NewTournament(name string) {
	c.cleanup()
	c.Tournament.SetName(name)
}

func (c *Context) cleanup() {
	c.Songs = NewUniqueList()
	c.Sound = NewSoundController(c.Log)
	c.Tournament = NewTournament("", c)
	c.Jingles = NewUniqueList()
	ChannelCleanup.Emit(EventTypeCleanup, struct{}{})
}

// StorageDir - return path to current tournament directory (and creates path if necessarry)
func (c *Context) StorageDir() string {
	p := path.Join(c.AppDir(), c.Tournament.Name)
	os.MkdirAll(p, 0700)
	return path.Join(p)
}

// MediaDir - return path to current tournament directory (and creates path if necessarry)
func (c *Context) MediaDir() string {
	p := path.Join(c.StorageDir(), "media")
	os.MkdirAll(p, 0700)
	return p
}

// AppDir - return path to application directory
func (c *Context) AppDir() string {
	u, _ := user.Current()
	p := path.Join(u.HomeDir, ".jinglemanager")
	os.MkdirAll(p, 0700)
	return path.Join(p)
}

// LastTournament - return path to application directory
func (c *Context) LastTournament() string {
	f, _ := os.Open(path.Join(c.AppDir(), "last.tournament"))
	t, _ := ioutil.ReadAll(f)
	return string(t)
}
