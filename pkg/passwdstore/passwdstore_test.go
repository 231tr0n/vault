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

	var tests = [][]string{
		{".vault", ".passwdstore"},
		{".vault", ".passwdstore"},
		{".vault", ".teststore"},
	}

	for _, test := range tests {
		t.Log(test)
		var passwdStoreFilePath = filepath.Join(append([]string{tempDir}, test...)...)

		var err = passwdstore.Init(passwdStoreFilePath)
		if err != nil {
			t.Fatal(err)
		}

		_, err = os.Stat(passwdStoreFilePath)
		if err != nil {
			failTestCase(t, passwdStoreFilePath, "File not found", "File should be created")
		}
	}
}

func TestChangeMasterPasswd(t *testing.T) {
	var tempDir = t.TempDir()

	var passwdStoreFilePath = filepath.Join(tempDir, ".vault", ".passwdstore")

	var tests = [][2][]byte{
		{
			[]byte("testsecret"),
			[]byte("secret"),
		},
	}

	var err = passwdstore.Init(passwdStoreFilePath)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Log(test)
		var err = passwdstore.Set("hi", "test", test[0])
		if err != nil {
			t.Fatal(err)
		}

		err = passwdstore.ChangeMasterPasswd(test[1], test[0])
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

func TestSetGetListDelete(t *testing.T) {
	var tempDir = t.TempDir()

	var passwdStoreFilePath = filepath.Join(tempDir, ".vault", ".passwdstore")

	var tests = [][3]string{
		{"hi", "how are you", "secret"},
	}

	var err = passwdstore.Init(passwdStoreFilePath)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Log(test)
		var passwd = []byte(test[2])

		var err = passwdstore.Set(test[0], test[1], passwd)
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

		var check = false
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
