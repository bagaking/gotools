package procast

import "sync"

type procCloseOrFailed struct {
	chErr  chan error
	chStop chan struct{}
	error  error
	closed bool
	mu     sync.Mutex
}

func (sc *procCloseOrFailed) C() <-chan struct{} {
	return sc.chStop
}

func (sc *procCloseOrFailed) Done() *procCloseOrFailed {
	sc.mu.Lock()
	if sc.closed == false {
		close(sc.chStop)
		sc.closed = true
	}
	sc.mu.Unlock()
	return sc
}

func (sc *procCloseOrFailed) Closed() bool {
	return sc.closed
}

func (sc *procCloseOrFailed) Err() error {
	return sc.error
}

func (sc *procCloseOrFailed) Fail(err error) {
	sc.chErr <- err
}

func NewCloseOrFailedProc(fn func(err error) error) (proc *procCloseOrFailed) {
	proc = &procCloseOrFailed{}
	proc.mu.Lock()
	defer proc.mu.Unlock()

	proc.chErr = make(chan error)
	proc.chStop = make(chan struct{})

	proc.Go(func() {
		for {
			select {
			case <-proc.C(): // finished
				return
			case err := <-proc.chErr:
				if fn == nil {
					proc.Done()
					continue
				}

				if err = fn(err); err != nil {
					proc.Done()
					return
				}
			}
		}
	})

	return proc
}

// GoAfterStop - panic of fn will not be recovered
func (sc *procCloseOrFailed) GoAfterStop(fn func(err error)) *procCloseOrFailed {
	go func() {
		<-sc.chStop
		fn(sc.error)
		return
	}()
	return sc
}

// Go - recover and close the chStop when fn panics
func (sc *procCloseOrFailed) Go(fn func()) *procCloseOrFailed {
	SafeGo(fn, func(err error) {
		sc.Done()
		sc.error = err
	})
	return sc
}