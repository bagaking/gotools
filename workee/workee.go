package workee

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bagaking/gotools/procast"
)

type (
	InitStrategy byte

	Workee interface {
		Name() string
		ID() string
		Close()
		IsFinished() bool
	}

	worker struct {
		id       string
		name     string
		finished bool

		conf

		chClose   chan struct{}
		watchOnce *sync.Once
	}
)

var _ Workee = &worker{}

const (
	InitAtLeastOnce InitStrategy = 0
	InitAsync       InitStrategy = 1

	StrWorkeeStart = "workee start"
	StrWorkeeRun   = "workee run"
	StrWorkeeExit  = "workee exit"
)

var DefaultConf = conf{
	TickDuration: time.Second,
	InitStrategy: InitAtLeastOnce,
	ErrorHandler: func(err error) {},
	ProcPrinter:  func(w Workee, str string, round int64) {},
}

func New(name string, fn func() error, opts ...Option) Workee {
	c := DefaultConf
	for _, opt := range opts {
		c = opt(c)
	}
	worker := &worker{
		conf:      c,
		chClose:   make(chan struct{}),
		watchOnce: &sync.Once{},
	}
	worker.name = name
	rand.Seed(time.Now().UnixNano())
	worker.id = fmt.Sprintf("%d-%d", rand.Uint64(), time.Now().UnixNano())
	worker.watchOnce.Do(func() {
		worker.start(fn)
	})
	return worker
}

func (w *worker) Name() string {
	return w.name
}

func (w *worker) ID() string {
	return w.id
}

func (w *worker) IsFinished() bool {
	return w.finished
}

func (w *worker) Close() {
	close(w.chClose)
}

func (w *worker) start(fn func() error) {
	w.ProcPrinter(w, StrWorkeeStart, 0)
	// once exit called, the holding will release
	// but the goroutine will go on until close
	_ = procast.HoldGo(func(exit func(err error)) {
		count := int64(0)

		if w.InitStrategy == InitAsync {
			exit(nil)
		} else {
			defer exit(nil)
		}

		defer procast.Recover(w.ErrorHandler, "!!! panic and quit")
		defer w.ProcPrinter(w, StrWorkeeExit, count)

		HoldAndTickUntilClose(w.TickDuration, func() {
			count++
			w.ProcPrinter(w, StrWorkeeRun, count)
			defer procast.Recover(w.ErrorHandler, "!!! panic")
			if err := fn(); err != nil && w.ErrorHandler != nil { // err will not lead the closing of ticks
				w.ErrorHandler(err)
			}
			exit(nil) // for the first return (a.k.a. atLeastOnce)
		}, w.chClose)

		w.finished = true
	})
}
