package passwdstore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/231tr0n/vault/pkg/passwdstore"
)

func failTestCase(t *testing.T, i any, o any, w any) {
	t.Error("Input:", i, "|", "Output:", o, "|", "Want:", w)
}

func TestInit(t *testing.T) {
	var tempDir = t.TempDir()

	var passwdStoreFilePath = filepath.Join(tempDir, ".vault", ".passwdstore")

	var err = passwdstore.Init(passwdStoreFilePath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(passwdStoreFilePath)
	if err != nil {
		failTestCase(t, passwdStoreFilePath, "File not found", "File should be created")
	}
}

func TestSet(t *testing.T) {
	var tempDir = t.TempDir()

	var passwdStoreFilePath = filepath.Join(tempDir, ".vault", ".passwdstore")

	var passwd = []byte("secret")

	var err = passwdstore.Init(passwdStoreFilePath)
	if err != nil {
		t.Fatal(err)
	}

	passwdstore.Set("hi", "how are you", passwd)
}
