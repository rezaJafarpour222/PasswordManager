package main

import (
	"fmt"
	"os"
	"pass/cli"
)

type Credential struct {
	username string
	password string
}

func main() {

	app := cli.New("0.0.8", "Vault.json")
	err := app.Run(os.Args, "HARDpassword1")
	if err != nil {
		fmt.Println(err)
	}

}
