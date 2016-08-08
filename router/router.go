package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/logger"
	"net/http"
)

// Router - Router to allow middlewares
type Router struct {
	*httprouter.Router
	*logger.Log
	middlewares []MiddlewareI
}

// MiddlewareI - interface needed to be implemented by middlewares
type MiddlewareI interface {
	ShouldPerform(http.ResponseWriter, *http.Request) bool
	Perform(http.ResponseWriter, *http.Request) error
	OnError(http.ResponseWriter, *http.Request)
}

// NewRouter - Will create new router
func NewRouter(l *logger.Log) *Router {
	return &Router{
		httprouter.New(),
		l,
		make([]MiddlewareI, 0),
	}
}

// AddMiddleware - will add mileware to perform on request
func (r *Router) AddMiddleware(m MiddlewareI) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// perform middleware actions here
	r.Log.Debug("You're hitting " + req.URL.String())
	for _, m := range r.middlewares {
		if m.ShouldPerform(w, req) {
			err := m.Perform(w, req)
			if err != nil {
				m.OnError(w, req)
				return
			}
		}
	}

	r.Router.ServeHTTP(w, req)
}
