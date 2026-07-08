package cli

import (
	"fmt"
	"pass/TUI"
	"strings"
)

func Confirmation(operation string, dp []TUI.DataPoint) bool {
	printer := TUI.Print{}
	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle(operation)

	box.PrintData(dp)
	printer.WithSecondary().PrintText("Are you sure? [y/n]")

	var input string
	fmt.Scan(&input)
	if strings.ToLower(input) == "y" {
		return true
	}
	printer.WithSecondary().PrintText(operation + " Canceled")
	fmt.Println()
	return false
}
