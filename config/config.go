package config

import (
	"os"
	"path/filepath"
)

func GetPasswdStoreFilePath() string {
	return filepath.Join(os.Getenv("HOME"), ".vault", ".passwdstore")
}
