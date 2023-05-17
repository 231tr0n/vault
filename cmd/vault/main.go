package main

import (
	"fmt"
	"os"

	"github.com/231tr0n/vault/internal/cli"
)

func main() {
	err := cli.Init()
	if err != nil {
		fmt.Println("-----------------")
		fmt.Println(err)
		fmt.Println("-----------------")
		os.Exit(1)
	}

	err = cli.Parse()
	if err != nil {
		fmt.Println("-----------------")
		fmt.Println(err)
		fmt.Println("-----------------")
		os.Exit(1)
	}
}
