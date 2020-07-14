package login

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// AESEncrypt encrypts plainData using key and returns the encrypted output.
// len(key) must be 16.
func AESEncrypt(key, plainData []byte) (nonceAndCipherData []byte, reterr error) {

	// just in case
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				reterr = fmt.Errorf("AESEncrypt caught panic: %w", err)
			} else {
				reterr = fmt.Errorf("AESEncrypt caught panic: %v", r)
			}
		}
	}()

	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(cb)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	cipherData := aead.Seal(nil, nonce, plainData, nil)
	return append(nonce, cipherData...), nil

}

// AESDecrypt reverses the action of AESEncrypt.
func AESDecrypt(key, nonceAndCipherData []byte) (plainData []byte, reterr error) {

	// just in case; iirc some decryption inputs can cause a panic
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				reterr = fmt.Errorf("AESDecrypt caught panic: %w", err)
			} else {
				reterr = fmt.Errorf("AESDecrypt caught panic: %v", r)
			}
		}
	}()

	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(cb)
	if err != nil {
		return nil, err
	}

	nonce := nonceAndCipherData[:aead.NonceSize()]
	cipherData := nonceAndCipherData[aead.NonceSize():]

	plainData, reterr = aead.Open(nil, nonce, cipherData, nil)
	return

}
