package datatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	var (
		err   error
		table = New(nil)
	)

	err = table.Set(1, 1, Plain("asd"))
	assert.Error(t, err, "cannot set out of range")
	err = table.Set(0, 1, Plain("asd"))
	assert.Error(t, err, "cannot set out of range")
	err = table.Set(1, 0, Plain("asd"))
	assert.Error(t, err, "cannot set out of range")
}
