package passwdstore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/231tr0n/vault/pkg/passwdstore"
)

func failTestCase(t *testing.T, i, o, w any) {
	t.Helper()
	t.Error("Input:", i, "|", "Output:", o, "|", "Want:", w)
}

func TestInit(t *testing.T) {
	tempDir := t.TempDir()

	tests := [][]string{
		{".vault", ".passwdstore"},
		{".vault", ".passwdstore"},
		{".vault", ".teststore"},
	}

	for _, test := range tests {
		t.Log(test)
		passwdStoreFilePath := filepath.Join(append([]string{tempDir}, test...)...)

		err := passwdstore.Init(passwdStoreFilePath)
		if err != nil {
			t.Fatal(err)
		}

		_, err = os.Stat(passwdStoreFilePath)
		if err != nil {
			failTestCase(t, passwdStoreFilePath, "File not found", "File should be created")
		}
	}
}

func TestChangePasswd(t *testing.T) {
	tempDir := t.TempDir()

	passwdStoreFilePath := filepath.Join(tempDir, ".vault", ".passwdstore")

	tests := [][2][]byte{
		{
			[]byte("testsecret"),
			[]byte("secret"),
		},
	}

	err := passwdstore.Init(passwdStoreFilePath)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Log(test)
		err := passwdstore.Set("hi", "test", test[0])
		if err != nil {
			t.Fatal(err)
		}

		err = passwdstore.ChangePasswd(test[1], test[0])
		if err != nil {
			t.Fatal(err)
		}

		var value string
		value, err = passwdstore.Get("hi", test[1])

		if err != nil {
			if err.Error() != "Wrong Password" {
				t.Fatal(err)
			}

			failTestCase(t, test, test[0], test[1])
		}

		if value != "test" {
			failTestCase(t, test, test[0], test[1])
		}
	}
}

func TestStoreMethods(t *testing.T) {
	tempDir := t.TempDir()

	passwdStoreFilePath := filepath.Join(tempDir, ".vault", ".passwdstore")

	tests := [][3]string{
		{"hi", "how are you", "secret"},
	}

	err := passwdstore.Init(passwdStoreFilePath)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Log(test)
		passwd := []byte(test[2])

		err := passwdstore.Set(test[0], test[1], passwd)
		if err != nil {
			t.Fatal(err)
		}

		var value string
		value, err = passwdstore.Get(test[0], passwd)

		if err != nil {
			t.Fatal(err)
		}

		if value != test[1] {
			failTestCase(t, test, value, test[1])
		}

		var list []string
		list, err = passwdstore.List(passwd)

		if err != nil {
			t.Fatal(err)
		}

		check := false

		for _, val := range list {
			if val == test[0] {
				check = true

				break
			}
		}

		if !check {
			failTestCase(t, test, "\"\"", test[0])
		}

		err = passwdstore.Delete(test[0], passwd)
		if err != nil {
			t.Fatal(err)
		}

		value, err = passwdstore.Get(test[0], passwd)
		if err != nil {
			t.Fatal(err)
		}

		if value != "" {
			failTestCase(t, test, value, test[1])
		}
	}
}
