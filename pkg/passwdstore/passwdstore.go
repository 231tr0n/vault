package passwdstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/231tr0n/vault/pkg/crypto"
	"github.com/231tr0n/vault/pkg/errorwrap"
)

const (
	fileComponents = 2
)

type passwdStore struct {
	Passwd []byte            `json:"passwd"`
	Store  map[string]string `json:"store"`
}

var (
	passwdStoreFilePath         = ""
	ErrFilePathNotAbsolute      = errors.New("given filepath not absolute")
	ErrPasswdFileManuallyEdited = errors.New("password file manually edited")
	ErrPasswdFileIntegrityFail  = errors.New("password file integrity fail")
	ErrVaultPasswdNotSet        = errors.New("vault password not set using ChangePasswd")
)

// Init sets the given filepath for password store file.
func Init(f string) error {
	if !filepath.IsAbs(f) {
		return errorwrap.ErrWrap(ErrFilePathNotAbsolute)
	}

	if stat, err := os.Stat(f); err == nil {
		if !stat.IsDir() {
			passwdStoreFilePath = f

			return nil
		}
	}

	err := os.MkdirAll(filepath.Dir(f), os.ModePerm)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	_, err = os.Create(f)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	passwdStoreFilePath = f

	return nil
}

// decryptFileData decrypts the file, unmarshals the json to struct and returns it.
func decryptFileData(p []byte) (passwdStore, error) {
	_, err := os.Stat(passwdStoreFilePath)
	if err != nil {
		return passwdStore{
			Passwd: []byte(""),
			Store:  make(map[string]string),
		}, errorwrap.ErrWrap(err)
	}

	data, err := os.ReadFile(passwdStoreFilePath)
	if err != nil {
		return passwdStore{
			Passwd: []byte(""),
			Store:  make(map[string]string),
		}, errorwrap.ErrWrap(err)
	}

	if len(data) == 0 {
		return passwdStore{
			Passwd: []byte(""),
			Store:  make(map[string]string),
		}, nil
	}

	pData := bytes.Split(data, []byte{'.'})
	if len(pData) != fileComponents {
		return passwdStore{
			Passwd: []byte(""),
			Store:  make(map[string]string),
		}, errorwrap.ErrWrap(ErrPasswdFileManuallyEdited)
	}

	h, err := crypto.Hash(pData[0], nil)
	if err != nil {
		return passwdStore{
			Passwd: []byte(""),
			Store:  make(map[string]string),
		}, errorwrap.ErrWrap(err)
	}

	if !crypto.Verify(h, pData[1]) {
		return passwdStore{
			Passwd: []byte(""),
			Store:  make(map[string]string),
		}, errorwrap.ErrWrap(ErrPasswdFileIntegrityFail)
	}

	s, err := crypto.Decrypt(pData[0], p)
	if err != nil {
		return passwdStore{
			Passwd: []byte(""),
			Store:  make(map[string]string),
		}, errorwrap.ErrWrap(err)
	}

	store := passwdStore{
		Passwd: []byte(""),
		Store:  make(map[string]string),
	}
	err = json.Unmarshal(s, &store)
	if err != nil {
		return passwdStore{
			Passwd: []byte(""),
			Store:  make(map[string]string),
		}, errorwrap.ErrWrap(err)
	}

	return store, nil
}

// encryptFileData marshals the struct to json, encrypts it and stores the content in the file.
func encryptFileData(store passwdStore, p []byte) error {
	if string(store.Passwd) == "" {
		return ErrVaultPasswdNotSet
	}

	store.Passwd = p
	s, err := json.Marshal(store)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	enc, err := crypto.Encrypt(s, p)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	h, err := crypto.Hash(enc, nil)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	err = os.WriteFile(passwdStoreFilePath, bytes.Join([][]byte{enc, h}, []byte{'.'}), os.ModePerm)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	return nil
}

// Get gets the key value pair from the store.
func Get(k string, p []byte) (string, error) {
	store, err := decryptFileData(p)
	if err != nil {
		return "", errorwrap.ErrWrap(err)
	}

	value := store.Store[k]

	return value, nil
}

// Put puts the key value pair in the store.
func Put(k, v string, p []byte) error {
	store, err := decryptFileData(p)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	store.Store[k] = v

	err = encryptFileData(store, p)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	return nil
}

// List lists all the key value pairs in the store.
func List(p []byte) ([]string, error) {
	store, err := decryptFileData(p)
	if err != nil {
		return nil, errorwrap.ErrWrap(err)
	}

	temp := make([]string, 0)
	for i := range store.Store {
		temp = append(temp, i)
	}

	return temp, nil
}

// Delete deletes the key value pair in the store.
func Delete(k string, p []byte) error {
	store, err := decryptFileData(p)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	delete(store.Store, k)

	err = encryptFileData(store, p)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	return nil
}

// Clear clears all the key value pairs in the store.
func Clear(p []byte) error {
	_, err := decryptFileData(p)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	empty := passwdStore{
		Passwd: []byte(""),
		Store:  make(map[string]string),
	}
	empty.Passwd = p
	err = encryptFileData(empty, p)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	return nil
}

// ChangePasswd changes the password for the store.
func ChangePasswd(np, op []byte) error {
	store, err := decryptFileData(op)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	store.Passwd = np

	err = encryptFileData(store, np)
	if err != nil {
		return errorwrap.ErrWrap(err)
	}

	return nil
}
