package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

const (
	aes256KeySize = 32
)

// ErrWrongPasswd is the error thrown when wrong password is given.
var ErrWrongPasswd = errors.New("crypto: wrong password")

func wrap(err error) error {
	if err != nil {
		return fmt.Errorf("crypto: %w", err)
	}

	return nil
}

// Encrypt encrypts "s" with password "p" using aes and gcm.
func Encrypt(s, p []byte) ([]byte, error) {
	if len(p)%aes256KeySize != 0 {
		temp := aes256KeySize - (len(p) % aes256KeySize)
		for i := 0; i < temp; i++ {
			p = append(p, '0')
		}
	}

	cr, err := aes.NewCipher(p)
	if err != nil {
		return nil, wrap(err)
	}

	gcm, err := cipher.NewGCM(cr)
	if err != nil {
		return nil, wrap(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, wrap(err)
	}

	return []byte(hex.EncodeToString(gcm.Seal(nonce, nonce, s, nil))), nil
}

// Decrypt decrypts "s" with password "p" using aes and gcm.
func Decrypt(s, p []byte) ([]byte, error) {
	s, err := hex.DecodeString(string(s))
	if err != nil {
		return nil, wrap(err)
	}

	if len(p)%aes256KeySize != 0 {
		temp := aes256KeySize - (len(p) % aes256KeySize)
		for i := 0; i < temp; i++ {
			p = append(p, '0')
		}
	}

	cr, err := aes.NewCipher(p)
	if err != nil {
		return nil, wrap(err)
	}

	gcm, err := cipher.NewGCM(cr)
	if err != nil {
		return nil, wrap(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ct := s[:nonceSize], s[nonceSize:]

	out, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, ErrWrongPasswd
	}

	return out, nil
}

// HmacHash hashes "s" with password "p" using hmac and sha256, appends it to "b" and returns it.
func HmacHash(s, p, b []byte) ([]byte, error) {
	hash := hmac.New(sha256.New, p)
	n, err := hash.Write(s)
	if n != len(s) || err != nil {
		return nil, wrap(err)
	}

	return []byte(hex.EncodeToString(hash.Sum(b))), nil
}

// Hash hashes "s" using sha256, appends it to "b" and returns it.
func Hash(s, b []byte) ([]byte, error) {
	hash := sha256.New()
	n, err := hash.Write(s)
	if n != len(s) || err != nil {
		return nil, wrap(err)
	}

	return []byte(hex.EncodeToString(hash.Sum(b))), nil
}

// Verify verifies if "s" is equal to "a".
func Verify(s, a []byte) bool {
	return string(s) == string(a)
}

// HmacVerify verifies if "s" is equal to "a" in a secure way.
// Use this for hmac based hashes.
func HmacVerify(s, a []byte) bool {
	return hmac.Equal(s, a)
}

// Generate generates a random byte array of length "s" and returns it.
// This function is used to generate random passwords.
func Generate(s int) ([]byte, error) {
	bytes := make([]byte, s)
	if _, err := rand.Read(bytes); err != nil {
		return nil, wrap(err)
	}

	return []byte(base64.StdEncoding.EncodeToString(bytes)[:s]), nil
}
