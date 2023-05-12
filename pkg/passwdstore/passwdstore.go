package passwdstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/231tr0n/vault/pkg/crypto"
)

var passwdStoreFilePath = ""

// Init sets the given filepath for password store file.
func Init(f string) error {
	if !filepath.IsAbs(f) {
		return errors.New("given filepath not absolute")
	}

	err := os.MkdirAll(filepath.Dir(f), 0o700)
	if err != nil {
		return err
	}

	_, err = os.Create(f)
	if err != nil {
		return err
	}

	passwdStoreFilePath = f

	return nil
}

func decryptFileData(p []byte) (map[string]string, error) {
	_, err := os.Stat(passwdStoreFilePath)
	if err != nil {
		return nil, err
	}

	var data []byte
	data, err = os.ReadFile(passwdStoreFilePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return make(map[string]string), nil
	}

	pData := bytes.Split(data, []byte{'-'})
	if len(pData) != 2 {
		return nil, errors.New("password file manually edited")
	}

	var h []byte
	h, err = crypto.Hash(pData[1], p, nil)
	if err != nil {
		return nil, err
	}

	if !crypto.Verify(h, pData[0]) {
		return nil, errors.New("password file integrity fail")
	}

	var s []byte
	s, err = crypto.Decrypt(pData[1], p)
	if err != nil {
		return nil, err
	}

	passwdStore := make(map[string]string)
	err = json.Unmarshal(s, &passwdStore)
	if err != nil {
		return nil, err
	}

	return passwdStore, nil
}

func encryptFileData(passwdStore map[string]string, p []byte) error {
	s, err := json.Marshal(passwdStore)
	if err != nil {
		return err
	}

	var enc []byte
	enc, err = crypto.Encrypt(s, p)
	if err != nil {
		return err
	}

	var h []byte
	h, err = crypto.Hash(enc, p, nil)
	if err != nil {
		return err
	}

	err = os.WriteFile(passwdStoreFilePath, bytes.Join([][]byte{h, enc}, []byte{'-'}), 0o700)
	if err != nil {
		return err
	}

	return nil
}

// Get gets the key value pair from the store.
func Get(k string, p []byte) (string, error) {
	passwdStore, err := decryptFileData(p)
	if err != nil {
		return "", err
	}

	value := passwdStore[k]
	return value, nil
}

// Set sets the key value pair in the store.
func Set(k string, v string, p []byte) error {
	passwdStore, err := decryptFileData(p)
	if err != nil {
		return err
	}

	passwdStore[k] = v

	err = encryptFileData(passwdStore, p)
	if err != nil {
		return err
	}

	return nil
}

// List lists all the key value pairs in the store.
func List(p []byte) ([]string, error) {
	passwdStore, err := decryptFileData(p)
	if err != nil {
		return nil, err
	}

	var temp []string
	for i := range passwdStore {
		temp = append(temp, i)
	}

	return temp, nil
}

// Delete deletes the key value pair provided in the store.
func Delete(k string, p []byte) error {
	passwdStore, err := decryptFileData(p)
	if err != nil {
		return err
	}

	delete(passwdStore, k)

	err = encryptFileData(passwdStore, p)
	if err != nil {
		return err
	}

	return nil
}

// ChangePasswd changes the password for the store.
func ChangePasswd(np []byte, op []byte) error {
	passwdStore, err := decryptFileData(op)
	if err != nil {
		return err
	}

	err = encryptFileData(passwdStore, np)
	if err != nil {
		return err
	}

	return nil
}
