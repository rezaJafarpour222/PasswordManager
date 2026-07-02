package cli

import (
	"fmt"
	"sync"
	"time"
)

func Spinner(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	frames := []rune{
		'⠋', '⠙', '⠹', '⠸',
		'⠼', '⠴', '⠦', '⠧',
		'⠇', '⠏',
	}
	i := 0

	for {
		select {
		case <-done:
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r%c Loading...", frames[i])
			i = (i + 1) % len(frames)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
