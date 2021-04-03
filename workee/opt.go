package workee

import "time"

type (
	conf struct {
		TickDuration time.Duration
		InitStrategy InitStrategy
		ErrorHandler func(err error)
		ProcPrinter  func(w *Worker, str string, round int64)
	}

	Option func(c conf) conf
)

func WithTickDuration(c conf, duration time.Duration) conf {
	c.TickDuration = duration
	return c
}

func WithInitStrategy(c conf, strategy InitStrategy) conf {
	c.InitStrategy = strategy
	return c
}

func WithErrorHandler(c conf, fn func(err error)) conf {
	c.ErrorHandler = fn
	return c
}

func WithProcPrinter(c conf, fn func(w *Worker, str string, round int64)) conf {
	c.ProcPrinter = fn
	return c
}
