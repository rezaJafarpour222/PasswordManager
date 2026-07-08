package cli

import (
	"fmt"
	"pass/TUI"
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

	printer := TUI.Print{}
	for {
		select {
		case <-done:
			fmt.Print("\r\033[K")
			return
		default:
			// fmt.Printf("\r%s%c%s %s%s%s",
			// 	TUI.Colors["Cyan"].Foreground(),
			// 	frames[i],
			// 	TUI.Reset,
			// 	TUI.Colors["Turquoise"].Foreground(),
			// 	spinnerText,
			// 	TUI.Reset,
			// )
			printer.WithAccent().PrintText("\r" + string(frames[i]))
			printer.WithPrimary().PrintText(spinnerText)
			i = (i + 1) % len(frames)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
