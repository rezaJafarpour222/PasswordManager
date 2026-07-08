package TUI

import (
	"fmt"
)

type DataPoint struct {
	Key   string
	Value string
}
type Box struct {
	Size          int
	LeftUpRune    rune
	RightUpRune   rune
	LeftDownRune  rune
	RightDownRune rune
	printer       Print
}

func NewBox(size int, leftUp, rightUp, leftDown, rightDown rune) *Box {
	return &Box{
		Size:          size,
		LeftUpRune:    leftUp,
		RightUpRune:   rightUp,
		LeftDownRune:  leftDown,
		RightDownRune: rightDown,
		printer:       Print{},
	}
}
func (b *Box) SetTitle(str string) {
	fmt.Println()
	b.printer.WithSecondary().PrintText(b.printer.Generate(" ", b.Size/2-len(str)) + str + "\n")
}

func (b *Box) PrintData(data []DataPoint) {
	b.printer.WithPrimary().PrintText(string(b.LeftUpRune) + b.printer.Generate("─", b.Size) + string(b.RightUpRune))
	fmt.Println()
	for _, field := range data {
		text := field.Key
		b.printer.WithSecondary().PrintText(text + ": ")
		b.printer.WithAccent().PrintText(
			field.Value + b.printer.Generate(" ", b.Size-(len(text)+len(field.Value))-1),
		)
		b.printer.WithPrimary().PrintText("│")
		fmt.Println()
	}
	b.printer.WithPrimary().PrintText(string(b.LeftDownRune) + b.printer.Generate("─", b.Size) + string(b.RightDownRune))
	fmt.Println()
}
