package TUI

import (
	"fmt"
	"strings"
)

type Print struct{}

func (p *Print) PrintText(str string) {
	fmt.Print(str)
	fmt.Print(Reset)

}

func (p *Print) WithPrimary() *Print {

	fmt.Print(Catppuccin["Lavender"].Foreground())

	return p
}

func (p *Print) WithSecondary() *Print {
	fmt.Print(Catppuccin["Mauve"].Foreground())

	return p
}

func (p *Print) WithAccent() *Print {
	fmt.Print(Catppuccin["Peach"].Foreground())
	return p
}

func (p *Print) WithDanger() *Print {
	fmt.Print(Catppuccin["Crust"].Foreground())
	return p
}

func (p *Print) Generate(str string, size int) string {
	return strings.Repeat(str, size)
}
