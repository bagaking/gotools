package lane

import (
	"context"
	"errors"
)

type contextKey struct {
	Name string
}

var (
	ErrLaneNotExist     = errors.New("target lane are not exist")
	ErrPayloadTypeError = errors.New("lane payload type error")
)

func (l *lane) ctxKey() contextKey {
	return contextKey{l.Name}
}

func (l *lane) TransportKey() string {
	return "__lane-" + l.Name
}

func (l *lane) CreateContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, l.ctxKey(), l.Payload)
	return ctx
}

func (l *lane) ExtractContext(ctx context.Context) error {
	val := ctx.Value(l.ctxKey())
	if val == nil {
		return ErrLaneNotExist
	}
	p, ok := val.(*Payloads)
	if !ok {
		return ErrPayloadTypeError
	}
	l.Payload = p

	return nil
}
