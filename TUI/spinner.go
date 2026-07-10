package TUI

import (
	"fmt"
	"sync"
	"time"
)

func Spinner(done <-chan struct{}, wg *sync.WaitGroup, spinnerText string) {
	defer wg.Done()

	frames := []rune{
		'⠋', '⠙', '⠹', '⠸',
		'⠼', '⠴', '⠦', '⠧',
		'⠇', '⠏',
	}
	i := 0

	printer := Print{}
	for {
		select {
		case <-done:
			fmt.Print("\r\033[K")
			return
		default:
			printer.WithAccent().PrintText("\r" + string(frames[i]))
			printer.WithPrimary().PrintText(spinnerText)
			i = (i + 1) % len(frames)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
