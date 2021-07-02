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
	type XX struct {
		V1 string `csv:"col=0"`
		V2 int    `csv:"col=1"`
	}

	type C struct {
		XX
		V3 float32 `csv:"col=2"`
		V4 string  `csv:"col=3"`
	}
	ret := make([]*C, 0, 2)

	reader := csv.NewReader(strings.NewReader(csvStr))

	err := ParseByCol(&ret, reader)

	assert.Nil(t, err)
	assert.Equal(t, *ret[0], C{XX{"a", 100}, 10e2, "s 1"})
	assert.Equal(t, *ret[1], C{XX{"b", 200}, 0.05, "s 2"})
}

func TestParseLineByColReturnsErrorForShortLine(t *testing.T) {
	type row struct {
		Name string `csv:"col=1"`
	}

	err := ParseLineByCol(&row{}, []string{"only-column"})
	if err == nil {
		t.Fatal("ParseLineByCol(&row{}, short line) error = nil, want non-nil")
	}
}

func TestParseLineByColCachesAnnotationsByStructType(t *testing.T) {
	type firstRow struct {
		Name string `csv:"col=0"`
	}
	type secondRow struct {
		Name string `csv:"col=1"`
	}

	first := firstRow{}
	if err := ParseLineByCol(&first, []string{"first", "unused"}); err != nil {
		t.Fatalf("ParseLineByCol(&firstRow{}, valid line) error = %v, want nil", err)
	}

	second := secondRow{}
	if err := ParseLineByCol(&second, []string{"unused", "second"}); err != nil {
		t.Fatalf("ParseLineByCol(&secondRow{}, valid line) error = %v, want nil", err)
	}
	if second.Name != "second" {
		t.Errorf("ParseLineByCol(&secondRow{}, valid line) set Name = %q, want %q", second.Name, "second")
	}
}
