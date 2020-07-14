// Package login provides simple username and password login functionality.
package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strings"
)

// ErrAuthnFailed is used to indicate login failure due to authentication (as opposed to other internal reasons).
var ErrAuthnFailed = errors.New("authentication failed")

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

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.URL.Query().Get("method"))
	if method == "" {
		method = r.Method
	}
	switch method {
	case "POST": // login
		h.postLogin(w, r)
	case "DELETE": // logout
		h.deleteLogin(w, r)
	default:
		http.Error(w, "Method not allowed", 405)
	}
	//case "GET": // get current login key - NOTE: this requires read, which we currently don't have
}

func (h *Handler) postLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type")); r.Method == "POST" && mt == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON parse error: %v", err), 400)
		}
	} else {
		in.Username = r.FormValue("username")
		in.Password = r.FormValue("password")
	}

	loginKey, err := h.lc.LoginCheck(in.Username, in.Password)
	if err != nil {
		if errors.Is(err, ErrAuthnFailed) {
			http.Error(w, "Authentication failed", 403)
			return
		}
		http.Error(w, "Unable to check login", 500)
		// FIXME: should have default logging mechanism plus a way to override on the handler
		log.Printf("login.Handler.postLogin: Unable to check login: %v", err)
		return
	}

	err = h.lw.LoginWrite(w, r, loginKey)
	if err != nil {
		http.Error(w, "Unable to write login", 500)
		// FIXME: should have default logging mechanism plus a way to override on the handler
		log.Printf("login.Handler.postLogin: LoginWrite error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var out struct {
		LoginKey string `json:"login_key"`
	}
	out.LoginKey = loginKey.LoginKey()
	_ = json.NewEncoder(w).Encode(&out)
}

func (h *Handler) deleteLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := h.lw.LoginWrite(w, r, nil)
	if err != nil {
		http.Error(w, "Unable to write login", 500)
		// FIXME: should have default logging mechanism plus a way to override on the handler
		log.Printf("login.Handler.deleteLogin: LoginWrite error: %v", err)
		return
	}

	w.WriteHeader(204)
}

// func (h *Handler) getLogin(w http.ResponseWriter, r *http.Request) {
// }

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
	// LoginRead returns who is currently logged in based on the current request.
	// If no login is provided then the Keyer returned must be nil.
	// A non-nil error can be returned to indicate this is an error condition,
	// but (nil,nil) is also a valid return and indicates nobody is logged in and
	// there was no error while determining this.
	LoginRead(w http.ResponseWriter, r *http.Request) (Keyer, error)
}

// ReaderFunc implements Reader as a function.
type ReaderFunc func(w http.ResponseWriter, r *http.Request) (Keyer, error)

// LoginRead implements Reader
func (f ReaderFunc) LoginRead(w http.ResponseWriter, r *http.Request) (Keyer, error) {
	return f(w, r)
}

// Writer knows how to write a login key to be associated with a request.
type Writer interface {
	// LoginWrite persists the login information from k so it is associated with this http response.
	// If k is non-nil then k.LoginKey() must not be empty.  Passing a nil k means the login info
	// should be removed (i.e. logout).
	LoginWrite(w http.ResponseWriter, r *http.Request, k Keyer) error
}

// WriterFunc implements Writer as a function.
type WriterFunc func(w http.ResponseWriter, r *http.Request, k Keyer) error

// LoginWrite implements Writer
func (f WriterFunc) LoginWrite(w http.ResponseWriter, r *http.Request, k Keyer) error {
	return f(w, r, k)
}
