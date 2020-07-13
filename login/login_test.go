package login

import (
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {

}

func ExampleHandler() {

	// import "github.com/d0sbit/authn/login"

	h := /*login.*/ New(nil, nil)
	fmt.Printf("HERE: %p", h)

	// Output: blah
}
