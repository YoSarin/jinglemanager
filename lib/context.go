package lib

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

func (c *Context) AppClosed() {
	c.Sound.ReleaseApps()
	c.Log.Info("Apps released")
	c.Save()
	c.Log.Close()
}
