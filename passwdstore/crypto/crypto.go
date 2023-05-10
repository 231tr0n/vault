package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

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

func Hash(s []byte, p []byte, b []byte) ([]byte, error) {
	var hash = hmac.New(sha256.New, p)
	var n, err = hash.Write(s)
	if n != len(s) || err != nil {
		return nil, err
	}
	return []byte(hex.EncodeToString(hash.Sum(b))), nil
}

func Verify(s []byte, a []byte) bool {
	return hmac.Equal(s, a)
}
