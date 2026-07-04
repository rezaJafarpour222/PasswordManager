package main

import (
	"fmt"
	"os"
	"pass/cli"
)

const VERSION string = "0.0.9"

func main() {

	app := cli.New(VERSION, "Vault.vault", "master.key")
	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(err)
	}

}
