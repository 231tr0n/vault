package main

import (
	"fmt"

	"github.com/231tr0n/vault/internal/cli"
)

func main() {
	err := cli.Init()
	if err != nil {
		fmt.Println(err)
	}

	err = cli.Parse()
	if err != nil {
		fmt.Println(err)
	}
}
