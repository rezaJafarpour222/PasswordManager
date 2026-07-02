package cli

import (
	"fmt"
	"time"
)

func Spinner(done <-chan struct{}) {
	frames := []rune{
		'⠋', '⠙', '⠹', '⠸',
		'⠼', '⠴', '⠦', '⠧',
		'⠇', '⠏',
	}
	i := 0

	for {
		select {
		case <-done:
			fmt.Print("\r \r") // Clear spinner
			return
		default:
			fmt.Printf("\r%c Loading...", frames[i])
			i = (i + 1) % len(frames)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
