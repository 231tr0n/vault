package passwdstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"github.com/231tr0n/vault/config"
	"github.com/231tr0n/vault/pkg/crypto"
)

func decryptFileData(p []byte) (map[string]string, error) {
	var _, err = os.Stat(config.PasswdStoreFile)
	if err != nil {
		return nil, err
	}

	var data []byte
	data, err = os.ReadFile(config.PasswdStoreFile)
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

	var passwordStore = make(map[string]string)
	err = json.Unmarshal(s, &passwordStore)
	if err != nil {
		return nil, err
	}

	return passwordStore, nil
}

func encryptFileData(passwordStore map[string]string, p []byte) error {
	var s, err = json.Marshal(passwordStore)
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

	err = os.WriteFile(config.PasswdStoreFile, bytes.Join([][]byte{h, enc}, []byte{'-'}), 0700)
	if err != nil {
		return err
	}

	return nil
}

func Get(k string, p []byte) (string, error) {
	var passwordStore, err = decryptFileData(p)
	if err != nil {
		return "", err
	}

	var value, ok = passwordStore[k]
	if ok {
		return value, nil
	}
	return "", errors.New("Not Found")
}

func Set(k string, v string, p []byte) error {
	var passwordStore, err = decryptFileData(p)
	if err != nil {
		return err
	}

	passwordStore[k] = v

	err = encryptFileData(passwordStore, p)
	if err != nil {
		return err
	}

	return nil
}

func List(p []byte) ([]string, error) {
	var passwordStore, err = decryptFileData(p)
	if err != nil {
		return nil, err
	}

	var temp []string
	for i := range passwordStore {
		temp = append(temp, i)
	}

	return temp, nil
}

func Delete(k string, p []byte) error {
	var passwordStore, err = decryptFileData(p)
	if err != nil {
		return err
	}

	delete(passwordStore, k)

	err = encryptFileData(passwordStore, p)
	if err != nil {
		return err
	}

	return nil
}

func ChangeMasterPassword(p []byte, op []byte) error {
	var passwordStore, err = decryptFileData(op)
	if err != nil {
		return err
	}

	err = encryptFileData(passwordStore, p)
	if err != nil {
		return err
	}

	return nil
}
