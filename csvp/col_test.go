package csvp

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	csvStr := `a,100,10e2,"s 1"
b,200,.05,"s 2"
`
	type C struct {
		V1 string  `csv:"col=0"`
		V2 int     `csv:"col=1"`
		V3 float32 `csv:"col=2"`
		V4 string  `csv:"col=3"`
	}
	ret := make([]*C, 0, 2)

	reader := csv.NewReader(strings.NewReader(csvStr))

	err := ParseByCol(&ret, reader)

	assert.Nil(t, err)
	assert.Equal(t, *ret[0], C{"a", 100, 10e2, "s 1"})
	assert.Equal(t, *ret[1], C{"b", 200, 0.05, "s 2"})
}
