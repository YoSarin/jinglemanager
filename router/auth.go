package router

import (
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
)

// AuthMiddleware - middleware for authentication
type AuthMiddleware struct {
	Log lib.LogI
}

// NewAuthMiddleware - will create new auth middleware
func NewAuthMiddleware(l lib.LogI) *AuthMiddleware {
	return &AuthMiddleware{
		Log: l,
	}
}

// ShouldPerform - check if middleware should run for this request
func (m *AuthMiddleware) ShouldPerform(*http.Request) bool {
	return true
}

// Perform - will perform middleware action
func (m *AuthMiddleware) Perform(*http.Request) error {
	m.Log.Debug("Auth performed")
	return nil
}

// OnError - action to be performed when there is an error
func (m *AuthMiddleware) OnError(http.ResponseWriter, *http.Request) {

}
