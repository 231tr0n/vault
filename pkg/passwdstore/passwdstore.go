package passwdstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/231tr0n/vault/pkg/crypto"
)

const (
	fileComponents = 2
)

type passwdStore struct {
	Passwd []byte            `json:"passwd"`
	Store  map[string]string `json:"store"`
}

func wrap(err error) error {
	if err != nil {
		return fmt.Errorf("passwdstore: %w", err)
	}

	return nil
}

func newpasswdStore() passwdStore {
	return passwdStore{
		Passwd: []byte(""),
		Store:  make(map[string]string),
	}
}

var (
	passwdStoreFilePath = ""
	// ErrFilePathNotAbsolute is the error thrown when
	// the passwdstore.Init function does not get an absolute path as an argument.
	ErrFilePathNotAbsolute = errors.New("passwdstore: given filepath not absolute")
	// ErrPasswdFileManuallyEdited is the error thrown when
	// the passwdstore file at $HOME/.vault/.passwdstore is not parsable as it is manually edited.
	ErrPasswdFileManuallyEdited = errors.New("passwdstore: password file manually edited")
	// ErrPasswdFileIntegrityFail is the error thrown when
	// the passwdstore file at $HOME/.vault/.passwdstore fails
	// the integrity check which means that the contents are edited.
	ErrPasswdFileIntegrityFail = errors.New("passwdstore: password file integrity fail")
	// ErrVaultPasswdNotSet is the error thrown when the
	// vault password is empty or not set. Use the passwdstore.ChangePasswd
	// function to set the password initially as it is empty.
	ErrVaultPasswdNotSet = errors.New("passwdstore: vault password not set")
)

// Init sets the given filepath for password store file.
func Init(f string) error {
	if !filepath.IsAbs(f) {
		return ErrFilePathNotAbsolute
	}

	if stat, err := os.Stat(f); err == nil {
		if !stat.IsDir() {
			passwdStoreFilePath = f

			return nil
		}
	}

	err := os.MkdirAll(filepath.Dir(f), os.ModePerm)
	if err != nil {
		return wrap(err)
	}

	_, err = os.Create(f)
	if err != nil {
		return wrap(err)
	}

	passwdStoreFilePath = f

	return nil
}

// decryptFileData decrypts the file, unmarshals the json to struct and returns it.
func decryptFileData(p []byte) (passwdStore, error) {
	_, err := os.Stat(passwdStoreFilePath)
	if err != nil {
		return newpasswdStore(), wrap(err)
	}

	data, err := os.ReadFile(passwdStoreFilePath)
	if err != nil {
		return newpasswdStore(), wrap(err)
	}

	if len(data) == 0 {
		return newpasswdStore(), nil
	}

	pData := bytes.Split(data, []byte{'.'})
	if len(pData) != fileComponents {
		return newpasswdStore(), ErrPasswdFileManuallyEdited
	}

	h, err := crypto.Hash(pData[0], nil)
	if err != nil {
		return newpasswdStore(), wrap(err)
	}

	if !crypto.Verify(h, pData[1]) {
		return newpasswdStore(), ErrPasswdFileIntegrityFail
	}

	s, err := crypto.Decrypt(pData[0], p)
	if err != nil {
		return newpasswdStore(), wrap(err)
	}

	store := newpasswdStore()
	err = json.Unmarshal(s, &store)
	if err != nil {
		return newpasswdStore(), wrap(err)
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
		return wrap(err)
	}

	enc, err := crypto.Encrypt(s, p)
	if err != nil {
		return wrap(err)
	}

	h, err := crypto.Hash(enc, nil)
	if err != nil {
		return wrap(err)
	}

	err = os.WriteFile(passwdStoreFilePath, bytes.Join([][]byte{enc, h}, []byte{'.'}), os.ModePerm)
	if err != nil {
		return wrap(err)
	}

	return nil
}

// Get gets the key value pair from the store.
func Get(k string, p []byte) (string, error) {
	store, err := decryptFileData(p)
	if err != nil {
		return "", wrap(err)
	}

	value := store.Store[k]

	return value, nil
}

// Put puts the key value pair in the store.
func Put(k, v string, p []byte) error {
	store, err := decryptFileData(p)
	if err != nil {
		return wrap(err)
	}

	store.Store[k] = v

	err = encryptFileData(store, p)
	if err != nil {
		return wrap(err)
	}

	return nil
}

// ListKeys lists all the keys in the store.
func ListKeys(p []byte) ([]string, error) {
	store, err := decryptFileData(p)
	if err != nil {
		return nil, wrap(err)
	}

	temp := make([]string, 0)
	for i := range store.Store {
		temp = append(temp, i)
	}

	return temp, nil
}

// ListEntries lists all the key value pairs in the store.
func ListEntries(p []byte) ([][2]string, error) {
	store, err := decryptFileData(p)
	if err != nil {
		return nil, wrap(err)
	}

	temp := make([][2]string, 0)
	for i, j := range store.Store {
		temp = append(temp, [2]string{i, j})
	}

	return temp, nil
}

// Delete deletes the key value pair in the store.
func Delete(k string, p []byte) error {
	store, err := decryptFileData(p)
	if err != nil {
		return wrap(err)
	}

	delete(store.Store, k)

	err = encryptFileData(store, p)
	if err != nil {
		return wrap(err)
	}

	return nil
}

// Clear clears all the key value pairs in the store.
func Clear(p []byte) error {
	_, err := decryptFileData(p)
	if err != nil {
		return wrap(err)
	}

	empty := newpasswdStore()
	empty.Passwd = p
	err = encryptFileData(empty, p)
	if err != nil {
		return wrap(err)
	}

	return nil
}

// ChangePasswd changes the password for the store.
func ChangePasswd(np, op []byte) error {
	store, err := decryptFileData(op)
	if err != nil {
		return wrap(err)
	}

	store.Passwd = np

	err = encryptFileData(store, np)
	if err != nil {
		return wrap(err)
	}

	return nil
}
