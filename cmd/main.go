package main

import (
	"fmt"
	"os"
	"pass/cli"
)

var (
	masterPassword string
)

const VERSION string = "0.0.8"

func main() {
	done := make(chan struct{})
	go cli.Spinner(done)
	app := cli.New(VERSION, "Vault.json")
	err := app.Run(os.Args, "HARDpassword1")
	if err != nil {
		fmt.Println(err)
	}

}
