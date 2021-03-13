package lane

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLane_Tag(t *testing.T) {
	var test Tag = "test"
	l := New("keepalive")
	assert.Equal(t, "", test.ReadFrom(l))
	test.WriteTo(l, "test_val")
	assert.Equal(t, "test_val", test.ReadFrom(l))
	test.ClearAt(l)
	assert.Equal(t, "", test.ReadFrom(l))
}

func TestLane_NewWithArgs(t *testing.T) {
	var test Tag = "_V_"
	l := New("keepalive", "ok")
	assert.Equal(t, "ok", Value.ReadFrom(l))
	test.WriteTo(l, "another_val")
	assert.Equal(t, "another_val", Value.ReadFrom(l))
}

func TestLane_Context(t *testing.T) {
	var test Tag = "ano"
	l := New("keepalive", "ok")
	test.WriteTo(l, "another_val")

	// this const should never change
	assert.Equal(t, "__lane-keepalive", l.TransportKey())

	ctx := l.CreateContext(context.Background())

	l2 := New("keepalive")
	err := l2.ExtractContext(ctx)

	assert.Nil(t, err)
	assert.Equal(t, "ok", Value.ReadFrom(l2))
	assert.Equal(t, "another_val", test.ReadFrom(l2))

	l2 = New("wrong_name")
	err = l2.ExtractContext(ctx)
	assert.ErrorIs(t, ErrLaneNotExist, err)

	l2 = New("keepalive")
	ctx = context.WithValue(ctx, contextKey{"keepalive"}, 1)
	err = l2.ExtractContext(ctx)
	assert.ErrorIs(t, ErrPayloadTypeError, err)
}

func TestLane_Payload(t *testing.T) {
	var test Tag = "ano"
	l := New("keepalive", "ok", test.String(), "v2")

	assert.True(t, Value.Of(l).Is("ok"), Value.ReadFrom(l))
	assert.True(t, test.Of(l).In("v1", "v2"))
	v, err := test.Of(l).Select(map[string]string{"v1": "1", "v2": "2"})
	assert.Nil(t, err)
	assert.Equal(t, "2", v)
	v, err = test.Of(l).Select(map[string]string{"v1": "1"})
	assert.ErrorIs(t, ErrCandidatesNotMatch, err)
	assert.Equal(t, "", v)
}
