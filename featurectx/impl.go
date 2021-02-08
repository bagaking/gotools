package featurectx

import (
	"context"
	"sync"
	"time"
)

type ctx struct {
	duration time.Duration

	createdAt time.Time
	startedAt time.Time
	closedAt  time.Time

	timerCtx    context.Context
	timerCancel func()

	chStart   chan struct{}
	chDone    chan struct{}
	onceStart sync.Once
}

var _ Context = &ctx{}

func (ctx *ctx) Start() (ok bool) {
	ctx.onceStart.Do(func() {
		ctx.startedAt = time.Now()
		ctx.timerCtx, ctx.timerCancel = context.WithTimeout(context.Background(), ctx.duration)
		close(ctx.chStart)
		ok = true
	})
	return
}

func (ctx *ctx) Abort() {
	ctx.Start()
	ctx.timerCancel()
	return
}

func (ctx *ctx) WaitForStart() {
	<-ctx.chStart
}

// Done cannot be called before MarkStart
func (ctx *ctx) WaitForDone() {
	<-ctx.Done()
}

func (ctx *ctx) Duration() time.Duration {
	return ctx.duration
}

func (ctx *ctx) Lasted() time.Duration {
	now := time.Now()
	if now.Before(ctx.startedAt) {
		return time.Duration(0)
	}
	return time.Now().Sub(ctx.startedAt)
}

func (ctx *ctx) Deadline() (deadline time.Time, ok bool) {
	if ctx.timerCtx == nil {
		return time.Now(), false
	}
	return ctx.timerCtx.Deadline()
}

func (ctx *ctx) Err() error {
	return ctx.timerCtx.Err()
}

func (ctx *ctx) Value(key interface{}) interface{} {
	return ctx.Value(key)
}

func (ctx *ctx) Done() <-chan struct{} {
	return ctx.chDone
}

func (ctx *ctx) CreatedAt() time.Time {
	return ctx.createdAt
}

func (ctx *ctx) StartedAt() time.Time {
	return ctx.startedAt
}

func (ctx *ctx) ClosedAt() time.Time {
	return ctx.createdAt
}

// New method returns a new featurectx.Context
// You have the `startTimeout` to start the context, otherwise the context will be closed directly
// After Start is called, the context will end after the `duration`
func New(duration time.Duration, startTimeout time.Duration) (c Context) {
	ctxNew := &ctx{
		duration:  duration,
		createdAt: time.Now(),

		chStart:   make(chan struct{}),
		chDone:    make(chan struct{}),
		onceStart: sync.Once{},
	}

	go func() {
		defer func() {
			ctxNew.closedAt = time.Now()
			close(ctxNew.chDone)
		}()
		select {
		case <-time.After(startTimeout):
			return
		case <-ctxNew.chStart:
		}
		<-ctxNew.timerCtx.Done()
	}()

	return ctxNew
}
