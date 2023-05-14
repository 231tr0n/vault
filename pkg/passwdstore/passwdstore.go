package passwdstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/231tr0n/vault/pkg/crypto"
)

const (
	fileComponents = 2
)

var (
	passwdStoreFilePath         = ""
	ErrFilePathNotAbsolute      = errors.New("passwdstore: given filepath not absolute")
	ErrPasswdFileManuallyEdited = errors.New("passwdstore: password file manually edited")
	ErrPasswdFileIntegrityFail  = errors.New("passwdstore: password file integrity fail")
)

func errWrap(err error) error {
	return errors.New("passwdstore: " + err.Error())
}

// Init sets the given filepath for password store file.
func Init(f string) error {
	if !filepath.IsAbs(f) {
		return ErrFilePathNotAbsolute
	}

	if _, err := os.Stat(f); err == nil {
		return nil
	}

	err := os.MkdirAll(filepath.Dir(f), os.ModePerm)
	if err != nil {
		return errWrap(err)
	}

	_, err = os.Create(f)
	if err != nil {
		return errWrap(err)
	}

	passwdStoreFilePath = f

	return nil
}

func decryptFileData(p []byte) (map[string]string, error) {
	_, err := os.Stat(passwdStoreFilePath)
	if err != nil {
		return nil, errWrap(err)
	}

	data, err := os.ReadFile(passwdStoreFilePath)
	if err != nil {
		return nil, errWrap(err)
	}

	if len(data) == 0 {
		return make(map[string]string), nil
	}

	pData := bytes.Split(data, []byte{'-'})
	if len(pData) != fileComponents {
		return nil, ErrPasswdFileManuallyEdited
	}

	h, err := crypto.Hash(pData[1], nil)
	if err != nil {
		return nil, errWrap(err)
	}

	if !crypto.Verify(h, pData[0]) {
		return nil, ErrPasswdFileIntegrityFail
	}

	s, err := crypto.Decrypt(pData[1], p)
	if err != nil {
		return nil, errWrap(err)
	}

	passwdStore := make(map[string]string)
	err = json.Unmarshal(s, &passwdStore)
	if err != nil {
		return nil, errWrap(err)
	}

	return passwdStore, nil
}

func encryptFileData(passwdStore map[string]string, p []byte) error {
	s, err := json.Marshal(passwdStore)
	if err != nil {
		return errWrap(err)
	}

	enc, err := crypto.Encrypt(s, p)
	if err != nil {
		return errWrap(err)
	}

	h, err := crypto.Hash(enc, nil)
	if err != nil {
		return errWrap(err)
	}

	err = os.WriteFile(passwdStoreFilePath, bytes.Join([][]byte{h, enc}, []byte{'-'}), os.ModePerm)
	if err != nil {
		return errWrap(err)
	}

	return nil
}

// Get gets the key value pair from the store.
func Get(k string, p []byte) (string, error) {
	passwdStore, err := decryptFileData(p)
	if err != nil {
		return "", errWrap(err)
	}

	value := passwdStore[k]

	return value, nil
}

// Set sets the key value pair in the store.
func Set(k, v string, p []byte) error {
	passwdStore, err := decryptFileData(p)
	if err != nil {
		return errWrap(err)
	}

	passwdStore[k] = v

	err = encryptFileData(passwdStore, p)
	if err != nil {
		return errWrap(err)
	}

	return nil
}

// List lists all the key value pairs in the store.
func List(p []byte) ([]string, error) {
	passwdStore, err := decryptFileData(p)
	if err != nil {
		return nil, errWrap(err)
	}

	temp := make([]string, 0)
	for i := range passwdStore {
		temp = append(temp, i)
	}

	return temp, nil
}

// Delete deletes the key value pair provided in the store.
func Delete(k string, p []byte) error {
	passwdStore, err := decryptFileData(p)
	if err != nil {
		return errWrap(err)
	}

	delete(passwdStore, k)

	err = encryptFileData(passwdStore, p)
	if err != nil {
		return errWrap(err)
	}

	return nil
}

// ChangePasswd changes the password for the store.
func ChangePasswd(np, op []byte) error {
	passwdStore, err := decryptFileData(op)
	if err != nil {
		return errWrap(err)
	}

	err = encryptFileData(passwdStore, np)
	if err != nil {
		return errWrap(err)
	}

	return nil
}
