package main

import (
	"os"
	"pass/TUI"
	"pass/cli"
)

const VERSION string = "1.1.0"

func main() {

	app := cli.New(VERSION, "Vault.vault", "Master.key")
	err := app.Run(os.Args)

	if err != nil {

		box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Error ", Value: err.Error()})
		box.PrintError(dp)
	}

}
