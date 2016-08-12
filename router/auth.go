package router

import (
	"encoding/base64"
	"errors"
	"github.com/martin-reznik/jinglemanager/lib"
	"net"
	"net/http"
	"strings"
)

// AuthMiddleware - middleware for authentication
type AuthMiddleware struct {
	Log     lib.LogI
	Context *lib.Context
}

// NewAuthMiddleware - will create new auth middleware
func NewAuthMiddleware(l lib.LogI, ctx *lib.Context) *AuthMiddleware {
	return &AuthMiddleware{
		Log:     l,
		Context: ctx,
	}
}

// ShouldPerform - check if middleware should run for this request
func (m *AuthMiddleware) ShouldPerform(w http.ResponseWriter, r *http.Request) bool {
	return true
}

// Perform - will perform middleware action
func (m *AuthMiddleware) Perform(w http.ResponseWriter, r *http.Request) error {
	stringIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(stringIP)
	if !m.Context.Tournament.Public && !ip.IsLoopback() {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return errors.New("Not Authorized")
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return errors.New("Not Authorized")
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return errors.New("Not Authorized")
		}

		for username, password := range m.Context.Tournament.Authorization {
			if pair[0] == username && pair[1] == password {
				return nil
			}
		}
		http.Error(w, "Not authorized", 401)
		return errors.New("Not Authorized")
	}
	return nil
}

// OnError - action to be performed when there is an error
func (m *AuthMiddleware) OnError(w http.ResponseWriter, r *http.Request) {
}
