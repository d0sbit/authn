package login

import (
	"bytes"
	"crypto/rand"
	mrand "math/rand"
	"sync"
	"testing"
)

func TestAESCrypt(t *testing.T) {

	key := make([]byte, 16)
	rand.Read(key)

	plain := []byte("abc123blah1")

	ct, err := AESEncrypt(key, plain)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Encrypted data: %X", ct)

	pt, err := AESDecrypt(key, ct)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pt, plain) {
		t.Errorf("expected equal values but got %q and %q", pt, plain)
	}

}

func TestAESDecryptFuzz(t *testing.T) {

	// send various garbage into AESDecrypt and make sure it doesn't panic

	c := 20000     // count of runs per goroutine
	gn := 4        // goroutines to launch
	maxlen := 1024 // max length of input

	key := make([]byte, 16)
	rand.Read(key)

	var wg sync.WaitGroup
	wg.Add(gn)
	for h := 0; h < gn; h++ {
		go func() {
			defer wg.Done()
			b := make([]byte, 0, maxlen)
			for i := 0; i < c; i++ {
				l := mrand.Intn(maxlen + 1)
				b = b[:l]
				rand.Read(b)
				_, err := AESDecrypt(key, b)
				if err == nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		}()
	}
	wg.Wait()

}
