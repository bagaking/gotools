package workee

import (
	"time"
)

// HoldAndTickUntilClose will execute fn per each tick
// if the param `tick` set with a value that are less than
// time.Microsecond, the interval will be set to time.Microsecond
// panic can cause the proc exit, the recover logic should be
// handled inside the handler `fn`
func HoldAndTickUntilClose(tick time.Duration, fn func(), chClose <-chan struct{}) {
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

