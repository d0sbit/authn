package login

import "net/http"

type Handler struct {
	Provider
}

type Provider interface {
	TryLogin(username, password string) (AuthnIDer, error)
}

type AuthnIDer interface {
	AuthnID() string
}

// AuthnIDString implements AuthnIDer as a string.
type AuthnIDString string

// AuthnID implements AuthnIDer.
func (s AuthnIDString) AuthnID() string {
	return string(s)
}

type LoginReader interface {
	ReadLogin(r *http.Request, w http.ResponseWriter) (AuthnIDer, error)
}

type LoginWriter interface {
	WriteLogin(r *http.Request, w http.ResponseWriter, a AuthnIDer) error
}
