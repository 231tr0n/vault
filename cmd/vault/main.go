package main

import (
	"fmt"
	"os"

	"github.com/231tr0n/vault/internal/cli"
)

func main() {
	err := cli.Init()
	if err != nil {
		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println(err)
		//nolint
		fmt.Println("-----------------")
		os.Exit(1)
	}

	err = cli.Parse()
	if err != nil {
		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println(err)
		//nolint
		fmt.Println("-----------------")
		os.Exit(1)
	}
}
