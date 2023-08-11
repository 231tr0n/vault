package main

import (
	"fmt"
	"os"

	"github.com/231tr0n/vault/internal/cli"
)

func main() {
	if err := cli.Init(); err != nil {
		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println(err)
		//nolint
		fmt.Println("-----------------")
		os.Exit(0)
	}

	if err := cli.Parse(); err != nil {
		//nolint
		fmt.Println("-----------------")
		//nolint
		fmt.Println(err)
		//nolint
		fmt.Println("-----------------")
		os.Exit(0)
	}
}
