package cli

import (
	"fmt"
	"strings"
)

func Confirmation(operation string, rest map[string]string) bool {

	fmt.Println("#======================================================#")
	for i, str := range rest {
		text := strings.ToUpper(string(i[0])) + i[1:]
		fmt.Printf("%s : %s\n", text, str)
	}
	fmt.Printf("Operation %s\n", operation)
	fmt.Printf("Are you sure? [y/n]")
	var input string
	fmt.Scan(&input)
	if strings.ToLower(input) == "y" {
		return true
	}
	fmt.Printf("\n%s Canceled.", operation)
	return false
}
