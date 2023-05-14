package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/231tr0n/vault/config"
)

func failTestCase(t *testing.T, i, o, w any) {
	t.Helper()
	t.Error("Input:", i, "|", "Output:", o, "|", "Want:", w)
}

func TestGetPasswdStoreFilePath(t *testing.T) {
	test := filepath.Join(os.Getenv("HOME"), ".vault", ".passwdstore")
	if config.GetPasswdStoreFilePath() != test {
		failTestCase(t, "\"\"", config.GetPasswdStoreFilePath(), test)
	}
}
