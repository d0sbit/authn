package login

import (
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {

}

func ExampleHandler() {

	h := /*login.*/ New(nil, nil)
	fmt.Printf("HERE: %p", h)

	// Output: blah
}
