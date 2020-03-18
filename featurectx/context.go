package featurectx

import (
	"context"
	"time"
)

type Context interface {
	context.Context

	// Start can be called multi times
	// The first marker takes effect at the beginning, subsequent calls are redundant.
	//
	// The main purpose of this design is to allow specific case to trigger the start-mark
	// in advance, and if not triggered manually, it will also be triggered automatically
	// in trigger layer methods, such as the Bench method.
	//
	// At the same time, the logic of the trigger layer should ensure that before the case
	// starts, it is Marked by the start
	Start() (ok bool)

	// Abort is used to close a Context
	// If a ctx is not yet started, the start channel is triggered (but the start log is
	// not printed), and then cancel it
	// If a ctx is already started, it will try to cancel the ctx
	Abort()

	// WaitForDone will stock the process until the ctx started
	WaitForStart()

	// WaitForDone will stock the process until the ctx finished or aborted
	// If a ctx has not yet started, it will wait for the ctx to start first
	WaitForDone()

	// Duration returns the configuration of duration
	Duration() time.Duration

	// Lasted returns the duration from start to now
	Lasted() time.Duration

	CreatedAt() time.Time
	StartedAt() time.Time
	ClosedAt() time.Time
}
