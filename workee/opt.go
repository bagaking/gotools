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

func WithTickDuration(duration time.Duration) Option {
	return func(c conf) conf {
		c.TickDuration = duration
		return c
	}
}

func WithInitStrategy(strategy InitStrategy) Option {
	return func(c conf) conf {
		c.InitStrategy = strategy
		return c
	}
}

func WithErrorHandler(fn func(err error)) Option {
	return func(c conf) conf {
		c.ErrorHandler = fn
		return c
	}
}

func WithProcPrinter(fn func(w *Worker, str string, round int64)) Option {
	return func(c conf) conf {
		c.ProcPrinter = fn
		return c
	}
}
