package cli

import (
	"flag"
	"fmt"
	"syscall"

	"github.com/231tr0n/vault/config"
	"github.com/231tr0n/vault/internal/errorwrap"
	"github.com/231tr0n/vault/pkg/crypto"
	"github.com/231tr0n/vault/pkg/passwdstore"
	"golang.org/x/term"
)

// Init initlialises the passwdstore.
func Init() error {
	return errorwrap.Wrap(passwdstore.Init(config.GetPasswdStoreFilePath()))
}

func readSecureInput(c string) ([]byte, error) {
	//nolint
	fmt.Print(c)

	//nolint
	s, err := term.ReadPassword(int(syscall.Stdin))

	//nolint
	fmt.Println()

	return s, errorwrap.Wrap(err)
}

// Parse parses the command line arguments and runs the respective functions accordingly.
func Parse() error {
	change := flag.Bool("change", false, "Changes the vault password.")
	list := flag.Bool("list", false, "Lists all the password names in the vault.")
	listAll := flag.Bool("list-all", false, "Lists all the passwords with names(dangerous) in the vault.")
	clear := flag.Bool("clear", false, "Clears all the passwords in the vault.")
	get := flag.String("get", "", "Gets the password from the vault.")
	put := flag.String("put", "", "Puts the password in the vault.")
	del := flag.String("delete", "", "Deletes the password in the vault.")
	//nolint
	generate := flag.Int("generate", 0, "Generates a new random password of length given. If this flag is passed along with put flag, it generates a random password and stores that in the vault.")

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
			return errorwrap.Wrap(err)
		}
		err = passwdstore.Clear(pwd)
		if err != nil {
			return errorwrap.Wrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Vault cleared")

	case *change:
		oldPwd, err := readSecureInput("Enter old vault password: ")
		if err != nil {
			return errorwrap.Wrap(err)
		}

		newPwd, err := readSecureInput("Enter new vault password: ")
		if err != nil {
			return errorwrap.Wrap(err)
		}

		newPwdCheck, err := readSecureInput("Re-Enter new vault password: ")
		if err != nil {
			return errorwrap.Wrap(err)
		}

		if string(newPwd) != string(newPwdCheck) {
			//nolint
			fmt.Println("-----------------")
			//nolint
			fmt.Println("New vault passwords don't match")

			break switch1
		}

		err = passwdstore.ChangePasswd(newPwd, oldPwd)
		if err != nil {
			return errorwrap.Wrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Vault password changed")

	case *listAll:
		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.Wrap(err)
		}

		list, err := passwdstore.ListEntries(pwd)
		if err != nil {
			return errorwrap.Wrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("List of passwords")
		//nolint
		fmt.Println("-----------------")

		for i, val := range list {
			//nolint
			fmt.Println(i, val[0], val[1])
		}

	case *list:
		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.Wrap(err)
		}

		list, err := passwdstore.ListKeys(pwd)
		if err != nil {
			return errorwrap.Wrap(err)
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
			return errorwrap.Wrap(err)
		}

		value, err := passwdstore.Get(*get, pwd)
		if err != nil {
			return errorwrap.Wrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Password:", value)

	case *put != "":
		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.Wrap(err)
		}

		var value []byte

		if *generate > 0 {
			value, err = crypto.Generate(*generate)
			if err != nil {
				return errorwrap.Wrap(err)
			}
			//nolint
			fmt.Println("-----------------")
			//nolint
			fmt.Println("Generated password:", string(value))
		} else {
			value, err = readSecureInput("Enter password for '" + *put + "': ")
			if err != nil {
				return errorwrap.Wrap(err)
			}
		}

		err = passwdstore.Put(*put, string(value), pwd)
		if err != nil {
			return errorwrap.Wrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Password stored")

	case *del != "":
		pwd, err := readSecureInput("Enter vault password: ")
		if err != nil {
			return errorwrap.Wrap(err)
		}

		err = passwdstore.Delete(*del, pwd)
		if err != nil {
			return errorwrap.Wrap(err)
		}

		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println("Password deleted")

	case *generate > 0:
		pwd, err := crypto.Generate(*generate)
		if err != nil {
			return errorwrap.Wrap(err)
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
