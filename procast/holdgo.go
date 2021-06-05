package procast

import "sync"

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

// HoldGo - Hold the proc until closer are called or panic
// closer can be called multi-times
// the quiting of HoldGo dose not means that the fn are finished
func HoldGo(fn func(closer func(error))) (err error) {
	wait, exit := AtLeastOnce()
	stop := func(e error) {
		if e != nil {
			err = e
		}
		exit()
	}
	SafeGo(func() {
		fn(stop)
	}, func(err error) {
		stop(err)
	})
	wait() // hold go
	return
}
