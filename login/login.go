// Package login provides simple username and password login functionality.
package login

import "net/http"

// New returns a new initialized Handler.
func New(lc Checker, lw Writer) *Handler {
	return &Handler{
		lc: lc,
		lw: lw,
	}
}

// Handler implements login API calls.
type Handler struct {
	lc Checker
	lw Writer
}

// Checker can check logins to see if they are valid and if so return the key.
type Checker interface {
	LoginCheck(username, password string) (Keyer, error)
}

// CheckerFunc implements Checker as a function
type CheckerFunc func(username, password string) (Keyer, error)

// LoginCheck implements Checker
func (f CheckerFunc) LoginCheck(username, password string) (Keyer, error) {
	return f(username, password)
}

// Keyer has a LoginKey method which is the unique ID associated with a login.
type Keyer interface {
	LoginKey() string // unique ID associated with login
}

// KeyString implements Keyer as a string.
type KeyString string

// LoginKey implements Keyer.
func (s KeyString) LoginKey() string {
	return string(s)
}

// Reader knows how to read login information from a request.
type Reader interface {
	LoginRead(r *http.Request, w http.ResponseWriter) (Keyer, error)
}

// ReaderFunc implements Reader as a function.
type ReaderFunc func(r *http.Request, w http.ResponseWriter) (Keyer, error)

// LoginRead implements Reader
func (f ReaderFunc) LoginRead(r *http.Request, w http.ResponseWriter) (Keyer, error) {
	return f(r, w)
}

// Writer knows how to write a login key to be associated with a request.
type Writer interface {
	LoginWrite(r *http.Request, w http.ResponseWriter, k Keyer) error
}

// WriterFunc implements Writer as a function.
type WriterFunc func(r *http.Request, w http.ResponseWriter, k Keyer) error

// LoginWrite implements Writer
func (f WriterFunc) LoginWrite(r *http.Request, w http.ResponseWriter, k Keyer) error {
	return f(r, w, k)
}
