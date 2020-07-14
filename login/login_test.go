package login

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
)

// func TestLogin(t *testing.T) {
// }

func ExampleHandler() {

	// encryption key for cookies
	cookieEncryptionKey := make([]byte, 16)
	// for this example we use random data but in a real app you
	// would read in or set a static value (in whatever form is secure enough for your purposes)
	rand.Read(cookieEncryptionKey)

	loginWriter := WriterFunc(func(w http.ResponseWriter, r *http.Request, k Keyer) error {
		return EncodeSecureJSONCookie(w, nil, cookieEncryptionKey, k.LoginKey())
	})

	loginReader := ReaderFunc(func(w http.ResponseWriter, r *http.Request) (Keyer, error) {
		var loginKeyString KeyString
		if err := DecodeSecureJSONCookie(r, "", cookieEncryptionKey, &loginKeyString); err != nil {
			return nil, err
		}
		return loginKeyString, nil
	})

	loginChecker := CheckerFunc(func(username, password string) (Keyer, error) {
		if username == "you@example.com" && password == "testtest" {
			return KeyString("u123"), nil
		}
		return nil, fmt.Errorf("login failed")
	})

	// simple example using mux - login api endpoint and another endpoint that uses the logged-in data
	mux := http.NewServeMux()

	// wire up login handler to respond to login API calls
	loginHandler := /*login.*/ New(loginChecker, loginWriter)
	mux.Handle("/api/login", loginHandler)

	// example of another handler that reads login data
	mux.HandleFunc("/api/demo", func(w http.ResponseWriter, r *http.Request) {
		l, _ := loginReader.LoginRead(w, r)
		if l == nil {
			w.WriteHeader(403)
			return
		}
		fmt.Fprintf(w, "%s", l.LoginKey())
	})

	// everything below here is a mock (test) client
	cjar, _ := cookiejar.New(nil)
	baseURL, _ := url.Parse("https://example.com/")

	// client login
	r := httptest.NewRequest("POST", baseURL.String()+"api/login",
		strings.NewReader(`{"username":"you@example.com","password":"testtest"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	fmt.Println(res.StatusCode) // print status code, should be 200
	cjar.SetCookies(baseURL, res.Cookies())

	// now hit the test endpoint that uses the login, passing the cookies
	r = httptest.NewRequest("GET", baseURL.String()+"api/demo", nil)
	for _, c := range cjar.Cookies(baseURL) {
		r.AddCookie(c)
	}
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res = w.Result()
	fmt.Println(res.StatusCode) // print status code, should be 200
	b, _ := ioutil.ReadAll(res.Body)
	// and print body contents, should be "u123" (our example of a user's unique ID/login key)
	fmt.Println(string(b))

	// hit it again without cookie and print the status code
	r = httptest.NewRequest("GET", baseURL.String()+"api/demo", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res = w.Result()
	fmt.Println(res.StatusCode) // print status code, should be 401

	// Output: 200
	// 200
	// u123
	// 403
}
