package featurectx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bagaking/gotools/featurectx"
	"github.com/stretchr/testify/assert"
)

const (
	ADDITIONAL_TIMEOUT_WAIT_TIME = time.Millisecond * 123
	MAX_DELTA_TIME               = time.Millisecond * 15
)

type ds []time.Duration

func (d ds) run(t *testing.T, f func(*testing.T, time.Duration)) bool {
	for _, duration := range durations {
		if result := t.Run(fmt.Sprintf("%v", duration), func(ts *testing.T) {
			f(ts, duration)
		}); !result {
			return result
		}
	}
	return true
}

var durations = ds{time.Millisecond * 50, time.Millisecond * 100, time.Millisecond * 234, time.Millisecond * 567}

func TestDuration(t0 *testing.T) {
	durations.run(t0, func(t *testing.T, duration time.Duration) {
		ctx := featurectx.New(duration)
		ctx.Start()
		start := time.Now()
		select {
		case <-time.After(duration + ADDITIONAL_TIMEOUT_WAIT_TIME):
			assert.Fail(t, "time out")
		case <-ctx.Done():
		}

		assert.WithinDuration(t, start.Add(duration), time.Now(), MAX_DELTA_TIME, "time not match")
	})
}

func TestShouldNotDoneWithoutAStart(t0 *testing.T) {
	durations.run(t0, func(t *testing.T, duration time.Duration) {
		ctx := featurectx.New(duration)
		select {
		case <-time.After(duration + ADDITIONAL_TIMEOUT_WAIT_TIME):
		case <-ctx.Done():
			assert.Fail(t, "")
		}
	})
}

func TestDelayStart(t0 *testing.T) {
	durations.run(t0, func(t *testing.T, duration time.Duration) {
		ctx := featurectx.New(duration)
		start := time.Now()
		startDelay := time.Millisecond * 233
		go func() {
			<-time.After(startDelay)
			ctx.Start()
		}()
		select {
		case <-time.After(duration + startDelay + ADDITIONAL_TIMEOUT_WAIT_TIME):
			assert.Fail(t, "time out")
		case <-ctx.Done():
		}
		assert.WithinDuration(t, start.Add(duration+startDelay), time.Now(), MAX_DELTA_TIME, "time not match")
	})
}

func TestMultiStart(t0 *testing.T) {
	durations.run(t0, func(t *testing.T, duration time.Duration) {
		ctx := featurectx.New(duration)
		start := time.Now()
		ctx.Start()
		go func() {
			<-time.After(duration / 4)
			ctx.Start()
			<-time.After(duration / 3)
			ctx.Start()
		}()
		select {
		case <-time.After(duration + ADDITIONAL_TIMEOUT_WAIT_TIME):
			assert.Fail(t, "time out")
		case <-ctx.Done():
		}
		assert.WithinDuration(t, start.Add(duration), time.Now(), MAX_DELTA_TIME, "time not match")
	})
}

func TestAbort(t0 *testing.T) {
	durations.run(t0, func(t *testing.T, duration time.Duration) {
		ctx := featurectx.New(duration)
		start := time.Now()
		ctx.Start()
		go func() {
			<-time.After(duration / 3)
			ctx.Abort()
		}()
		select {
		case <-time.After(duration + ADDITIONAL_TIMEOUT_WAIT_TIME):
			assert.Fail(t, "time out")
		case <-ctx.Done():
		}
		assert.WithinDuration(t, start.Add(duration/3), time.Now(), MAX_DELTA_TIME, "time not match")
	})
}

func TestAbortNotStatedCtx(t0 *testing.T) {
	durations.run(t0, func(t *testing.T, duration time.Duration) {
		ctx := featurectx.New(duration)
		start := time.Now()
		go func() {
			<-time.After(duration)
			ctx.Abort()
		}()
		select {
		case <-time.After(duration + ADDITIONAL_TIMEOUT_WAIT_TIME):
			assert.Fail(t, "time out")
		case <-ctx.Done():
		}
		assert.WithinDuration(t, start.Add(duration), time.Now(), MAX_DELTA_TIME, "time not match")
	})
}

func TestMultiAbort(t0 *testing.T) {
	durations.run(t0, func(t *testing.T, duration time.Duration) {
		ctx := featurectx.New(duration)
		start := time.Now()
		ctx.Start()
		go func() {
			<-time.After(duration / 3)
			ctx.Abort()
			<-time.After(duration / 2)
			ctx.Abort()
		}()
		select {
		case <-time.After(duration + ADDITIONAL_TIMEOUT_WAIT_TIME):
			assert.Fail(t, "time out")
		case <-ctx.Done():
		}
		assert.WithinDuration(t, start.Add(duration/3), time.Now(), MAX_DELTA_TIME, "time not match")
		<-time.After(duration)
	})
}
