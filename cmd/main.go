package main

import (
	"fmt"
	"os"
	"pass/cli"
)

const VERSION string = "1.0.0"

func main() {

	app := cli.New(VERSION, "Vault.vault", "Master.key")
	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(err)
	}

}
