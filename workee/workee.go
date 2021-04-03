package workee

import (
	"sync"
	"time"

	"github.com/bagaking/gotools/procast"
)

type (
	InitStrategy byte

	Worker struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		conf

		chClose   chan struct{}
		watchOnce *sync.Once
	}
)

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
	ProcPrinter:  func(w *Worker, str string, round int64) {},
}

func New(name string, fn func() error, opts ...Option) *Worker {
	c := DefaultConf
	for _, opt := range opts {
		c = opt(c)
	}
	worker := &Worker{
		conf:      c,
		chClose:   make(chan struct{}),
		watchOnce: &sync.Once{},
	}
	worker.Name = name
	worker.watchOnce.Do(func() {
		worker.start(fn)
	})
	return worker
}

func (w *Worker) Close() {
	close(w.chClose)
}

func (w *Worker) start(fn func() error) {
	wait, exit := AtLeastOnce()
	w.ProcPrinter(w, StrWorkeeStart, 0)

	go func() {
		count := int64(0)

		if w.InitStrategy == InitAsync {
			exit()
		} else {
			defer exit()
		}

		defer procast.Recover(w.ErrorHandler, "!!! panic and quit")
		defer w.ProcPrinter(w, StrWorkeeExit, count)

		UntilClose(w.TickDuration, func() {
			count++
			w.ProcPrinter(w, StrWorkeeRun, count)
			defer procast.Recover(w.ErrorHandler, "!!! panic")
			err := fn()
			exit()
			if w.ErrorHandler != nil {
				w.ErrorHandler(err)
			}
		}, w.chClose)
	}()
	wait()
}
