package cli

import (
	// Used by the commented code.
	_ "bufio"
	"flag"
	"fmt"
	// Used by the commented code.
	_ "os"
	"syscall"

	"github.com/231tr0n/vault/config"
	"github.com/231tr0n/vault/pkg/crypto"
	"github.com/231tr0n/vault/pkg/errorwrap"
	"github.com/231tr0n/vault/pkg/passwdstore"
	"golang.org/x/term"
)

// Init initlialises the passwdstore.
func Init() error {
	return errorwrap.ErrWrap(passwdstore.Init(config.GetPasswdStoreFilePath()))
}

func readSecureInput(c string) ([]byte, error) {
	//nolint
	fmt.Print(c)

	// fmt.Print("\033[?25l\033[8m")

	// stdinReader := bufio.NewReader(os.Stdin)
	// s, err := stdinReader.ReadBytes('\n')

	// fmt.Print("\033[28m\033[?25h")

	//nolint
	s, err := term.ReadPassword(int(syscall.Stdin))

	return s, errorwrap.ErrWrap(err)
}

// Parse parses the command line arguments and runs the respective functions accordingly.
func Parse() error {
	change := flag.Bool("change", false, "Changes the vault password")
	list := flag.Bool("list", false, "Lists all the passwords in the vault")
	clear := flag.Bool("clear", false, "Clears all the passwords in the vault")
	get := flag.String("get", "", "Gets the password from the vault")
	put := flag.String("put", "", "Puts the password in the vault")
	del := flag.String("delete", "", "Deletes the password in the vault")
	generate := flag.Int("generate", 0, "Generates a new random password of length given")

	flag.Parse()

	//nolint
	fmt.Println("-----------------")
	//nolint
	fmt.Println("Vault")
	//nolint
	fmt.Println("-----------------")

switch1:
	switch {
	case *clear:
		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}
		err = passwdstore.Clear(pwd)
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Vault cleared")

	case *change:
		oldPwd, err := readSecureInput("Enter old vault password: ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		newPwd, err := readSecureInput("Enter new vault password: ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		newPwdCheck, err := readSecureInput("Enter new vault password once again: ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		if string(newPwd) != string(newPwdCheck) {
			//nolint
			fmt.Println("-----------------")
			//nolint
			fmt.Println("Passwords don't match")

			break switch1
		}

		err = passwdstore.ChangePasswd(newPwd, oldPwd)
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Vault password changed")

	case *list:
		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		list, err := passwdstore.List(pwd)
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("List of passwords")
		//nolint
		fmt.Println("-----------------")

		for i, val := range list {
			//nolint
			fmt.Println(i, val)
		}

	case *get != "":
		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		value, err := passwdstore.Get(*get, pwd)
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Password:", value)

	case *put != "":
		value, err := readSecureInput("Enter password for " + *put + ": ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		err = passwdstore.Put(*put, string(value), pwd)
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Password stored")

	case *del != "":
		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		err = passwdstore.Delete(*del, pwd)
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Password deleted")

	case *generate > 0:
		pwd, err := crypto.Generate(*generate)
		if err != nil {
			return errorwrap.ErrWrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Generated password:", string(pwd))
	default:
		//nolint
		fmt.Println("No arguments given. Run -help to get a list of arguments.")
	}

	//nolint
	fmt.Println("-----------------")

	return nil
}
