package procast

import (
	"testing"
)

func TestRecover(t *testing.T) {
	func() {
		defer Recover(func(err error) {
			if err == nil {
				t.Error("recover failed")
			}
			if err.Error() != "recover, ok, x" {
				t.Error("recover content error")
			}
		}, "recover, %s", "ok")
		panic("x")
	}()
}
