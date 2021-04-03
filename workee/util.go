package workee

import (
	"sync"
	"time"
)

// UntilClose will execute fn per each tick
// if the param `tick` set with a value that are less than
// time.Microsecond, the interval will be set to time.Microsecond
// panic can cause the proc exit, the recover logic should be
// handled inside the handler `fn`
func UntilClose(tick time.Duration, fn func(), chClose <-chan struct{}) {
	var ticker *time.Ticker
	if tick > time.Microsecond {
		ticker = time.NewTicker(tick)
	} else {
		ticker = time.NewTicker(time.Microsecond)
	}

	for {
		fn() // error can be handled in the fn, do not panic or thrown

		select {
		case <-ticker.C:
		case <-chClose:
			ticker.Stop()
			return
		}
	}
}

// AtLeastOnce for returning a pair of functions.
// The `wait` method can block a goroutine, and the `exit` method causes
// this blocking to exit.
// `exit` can be called before or after the `wait`,
// It also can be called multiple times, but only the first call makes sense.
func AtLeastOnce() (wait func(), exit func()) {
	wg, once := sync.WaitGroup{}, sync.Once{}
	wg.Add(1)
	return wg.Wait, func() {
		once.Do(func() {
			wg.Done()
		})
	}
}
