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

// Init sets the given filepath for password store file
func Init(f string) error {
	if !filepath.IsAbs(f) {
		return errors.New("Given filepath not absolute.")
	}

	var err = os.MkdirAll(filepath.Dir(f), 0700)
	if err != nil {
		return err
	}

	_, err = os.Create(passwdStoreFilePath)
	if err != nil {
		return err
	}

	passwdStoreFilePath = f

	return nil
}

func decryptFileData(p []byte) (map[string]string, error) {
	var _, err = os.Stat(passwdStoreFilePath)
	if err != nil {
		return nil, err
	}

	var data []byte
	data, err = os.ReadFile(passwdStoreFilePath)
	if err != nil {
		return nil, err
	}

	var pData = bytes.Split(data, []byte{'-'})
	if len(data) == 0 {
		return nil, errors.New("No password data in file to decrypt file data")
	} else if len(pData) != 2 {
		return nil, errors.New("Password File Manually Edited")
	}

	var h []byte
	h, err = crypto.Hash(pData[1], p, nil)
	if err != nil {
		return nil, err
	}

	if crypto.Verify(h, pData[0]) {
		return nil, errors.New("Password File Manually Edited")
	}

	var s []byte
	s, err = crypto.Decrypt(pData[1], p)
	if err != nil {
		return nil, err
	}

	var passwdStore = make(map[string]string)
	err = json.Unmarshal(s, &passwdStore)
	if err != nil {
		return nil, err
	}

	return passwdStore, nil
}

func encryptFileData(passwdStore map[string]string, p []byte) error {
	var s, err = json.Marshal(passwdStore)
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

	err = os.WriteFile(passwdStoreFilePath, bytes.Join([][]byte{h, enc}, []byte{'-'}), 0700)
	if err != nil {
		return err
	}

	return nil
}

// Get gets the key value pair from the store
func Get(k string, p []byte) (string, error) {
	var passwdStore, err = decryptFileData(p)
	if err != nil {
		return "", err
	}

	var value, ok = passwdStore[k]
	if ok {
		return value, nil
	}
	return "", errors.New("Not Found")
}

// Set sets the key value pair in the store
func Set(k string, v string, p []byte) error {
	var passwdStore, err = decryptFileData(p)
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

// List lists all the key value pairs in the store
func List(p []byte) ([]string, error) {
	var passwdStore, err = decryptFileData(p)
	if err != nil {
		return nil, err
	}

	var temp []string
	for i := range passwdStore {
		temp = append(temp, i)
	}

	return temp, nil
}

// Delete deletes the key value pair provided in the store
func Delete(k string, p []byte) error {
	var passwdStore, err = decryptFileData(p)
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

// ChangeMasterPasswd changes the password for the store
func ChangeMasterPasswd(p []byte, op []byte) error {
	var passwdStore, err = decryptFileData(op)
	if err != nil {
		return err
	}

	err = encryptFileData(passwdStore, p)
	if err != nil {
		return err
	}

	return nil
}
