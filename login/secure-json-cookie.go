package login

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EncodeSecureJSONCookie is a convenience function that calls json.Marshal on obj, encrypts
// the result with AESEncrypt(key, ...) and sets it has a cookie on the HTTP response.
// TODO: describe defaults
func EncodeSecureJSONCookie(w http.ResponseWriter, c *http.Cookie, key []byte, obj interface{}) error {

	if c == nil {
		c = &http.Cookie{
			Name:     "login",
			MaxAge:   int((time.Hour * 8) / time.Second),
			HttpOnly: true,
		}
	}

	if c.Name == "" {
		c.Name = "login"
	}

	if c.MaxAge == 0 {
		c.MaxAge = int((time.Hour * 8) / time.Second)
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %w", err)
	}

	out, err := AESEncrypt(key, b)
	if err != nil {
		return err
	}

	c.Value = base64.RawURLEncoding.EncodeToString(out)

	http.SetCookie(w, c)

	return nil
}

// DecodeSecureJSONCookie reverses the action of EncodeSecureJSONCookie.
// The obj is what will be passed to json.Unmarshal (it should be a pointer).
// TODO: describe defaults, http.ErrNoCookie
func DecodeSecureJSONCookie(r *http.Request, cookieName string, key []byte, obj interface{}) error {

	if cookieName == "" {
		cookieName = "login"
	}

	c, err := r.Cookie(cookieName)
	if err != nil {
		return err // if no such cookie then http.ErrNoCookie will be returned
	}

	ct, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil {
		return err
	}
	pt, err := AESDecrypt(key, ct)
	if err != nil {
		return err
	}

	err = json.Unmarshal(pt, obj)
	if err != nil {
		return fmt.Errorf("json.Unmarshal error: %w", err)
	}

	return nil
}
