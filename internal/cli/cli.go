// Internal package to parse cli arguments and run respective functionality.
package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/231tr0n/vault/config"
	"github.com/231tr0n/vault/pkg/crypto"
	"github.com/231tr0n/vault/pkg/passwdstore"
)

func Init() error {
	return passwdstore.Init(config.GetPasswdStoreFilePath())
}

func readSecureInput(c string) ([]byte, error) {
	fmt.Print(c)
	fmt.Print("\033[?25l\033[8m")
	stdinReader := bufio.NewReader(os.Stdin)
	s, err := stdinReader.ReadBytes('\n')
	fmt.Print("\033[28m\033[?25h")
	return s, err
}

func Parse() error {
	change := flag.Bool("change", false, "Changes the vault password")
	list := flag.Bool("list", false, "Lists all the passwords in the vault")
	get := flag.String("get", "", "Gets the password from the vault")
	put := flag.String("put", "", "Puts the password with its identifier into the vault")
	del := flag.String("delete", "", "Deletes the password from the vault")
	generate := flag.Int("generate", 0, "Generates a new random password of length provided")
	flag.Parse()

	fmt.Println("-----------------")
	fmt.Println("Vault")
	fmt.Println("-----------------")

	if *change {
		oldPwd, err := readSecureInput("Enter old password:")
		if err != nil {
			return err
		}
		newPwd, err := readSecureInput("Enter new password:")
		if err != nil {
			return err
		}
		err = passwdstore.ChangePasswd(newPwd, oldPwd)
		if err != nil {
			return err
		}
		fmt.Println("-----------------")
		fmt.Println("Password Changed")
	} else if *list {
		pwd, err := readSecureInput("Enter vault password:")
		if err != nil {
			return err
		}
		list, err := passwdstore.List(pwd)
		if err != nil {
			return err
		}
		fmt.Println("-----------------")
		fmt.Println("List of passwords")
		fmt.Println("-----------------")
		fmt.Println(list)
		for _, val := range list {
			fmt.Println(val)
		}
	} else if *get != "" {
		pwd, err := readSecureInput("Enter vault password:")
		if err != nil {
			return err
		}
		value, err := passwdstore.Get(*get, pwd)
		if err != nil {
			return err
		}
		fmt.Println("-----------------")
		fmt.Println(value)
	} else if *put != "" {
		value, err := readSecureInput("Enter password to be stored:")
		if err != nil {
			return err
		}
		pwd, err := readSecureInput("Enter vault password:")
		if err != nil {
			return err
		}
		err = passwdstore.Set(*put, string(value), pwd)
		if err != nil {
			return err
		}
		fmt.Println("-----------------")
		fmt.Println("Password stored")
	} else if *del != "" {
		pwd, err := readSecureInput("Enter vault password:")
		if err != nil {
			return err
		}
		err = passwdstore.Delete(*del, pwd)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		fmt.Println("-----------------")
		fmt.Println("Password deleted")
	} else if *generate > 0 {
		pwd, err := crypto.Generate(*generate)
		if err != nil {
			return err
		}
		fmt.Println("-----------------")
		fmt.Println("Password generated:", pwd)
	} else {
		fmt.Println("No arguments given. Run -help to get a list of arguments.")
	}

	fmt.Println("-----------------")
	return nil
}
