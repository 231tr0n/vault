package config

import (
	"os"
	"path/filepath"
)

var (
	PasswdStoreFile = filepath.Join(os.Getenv("HOME"), ".vault", ".passwdstore")
)
