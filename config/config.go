package config

import (
	"os"
	"path/filepath"
)

var passwdStoreFilePath = filepath.Join(os.Getenv("HOME"), ".vault", ".passwdstore")

func GetPasswdStoreFilePath() string {
	return passwdStoreFilePath
}
