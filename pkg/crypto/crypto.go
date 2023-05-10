package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// Encryt encrypts "s" with password "p" using aes and gcm
func Encrypt(s []byte, p []byte) ([]byte, error) {
	if len(p)%32 != 0 {
		var temp = 32 - (len(p) % 32)
		for i := 0; i < temp; i++ {
			p = append(p, '0')
		}
	}

	var cr, err = aes.NewCipher(p)
	if err != nil {
		return nil, err
	}

	var gcm cipher.AEAD
	gcm, err = cipher.NewGCM(cr)
	if err != nil {
		return nil, err
	}

	var nonce = make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return []byte(hex.EncodeToString(gcm.Seal(nonce, nonce, s, nil))), nil
}

// Decrypt decrypts "s" with password "p" using aes and gcm
func Decrypt(s []byte, p []byte) ([]byte, error) {
	var err error
	s, err = hex.DecodeString(string(s))
	if err != nil {
		return nil, err
	}

	if len(p)%32 != 0 {
		var temp = 32 - (len(p) % 32)
		for i := 0; i < temp; i++ {
			p = append(p, '0')
		}
	}

	var cr cipher.Block
	cr, err = aes.NewCipher(p)
	if err != nil {
		return nil, err
	}

	var gcm cipher.AEAD
	gcm, err = cipher.NewGCM(cr)
	if err != nil {
		return nil, err
	}

	var nonceSize = gcm.NonceSize()
	var nonce, ct = s[:nonceSize], s[nonceSize:]

	var out []byte
	out, err = gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Hash hashes "s" with password "p" using hmac and sha256, appends it to "b" and returns it
func Hash(s []byte, p []byte, b []byte) ([]byte, error) {
	var hash = hmac.New(sha256.New, p)
	var n, err = hash.Write(s)
	if n != len(s) || err != nil {
		return nil, err
	}
	return []byte(hex.EncodeToString(hash.Sum(b))), nil
}

// Verify verifies if "s" is equal to "a" in a secure way
func Verify(s []byte, a []byte) bool {
	return hmac.Equal(s, a)
}

// Generate generates a random byte array of length "s" and returns it. This function is used to generate random passwords.
func Generate(s int) ([]byte, error) {
	var bytes = make([]byte, s)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(bytes)[:s]), nil
}
