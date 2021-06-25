package datatable

import (
	"encoding/json"
	"errors"
	"fmt"
)

type (
	Value fmt.Stringer

	Grid  []Line
	Plain string

	SingleQuery struct {
		Tag string `json:"tag"`
		Val Plain  `json:"val"`
	}

	Table interface {
		GetTitleLine() *TitleLine

		/* index */

		// Width - get width of the table
		Width() int
		Height() int
		MaxRow() int
		GetRowByColAndVal(col int, value Value) int

		/* values */

		// AppendRow - append a row in the end of the table
		AppendRow(line Line) error
		SetRow(row int, line Line) error
		GetRow(row int) Line

		Set(row, col int, val Value) error

		GetLine(row int) Line
		// Get returns Empty when out of range
		Get(row, col int) Value
		Gets(row int, cols ...int) Line

		/* high level */

		// FindOnePos returns the row and col by tag and tagVal
		FindOnePos(q SingleQuery) (row int, col int)

		// FindOne returns the value by tag and tagVal
		FindOne(q SingleQuery, selectTags ...string) Line

		Render(titles TitleLine, showTitle bool) [][]string

		json.Marshaler
		json.Unmarshaler
	}
)

const Empty Plain = ""

func (t Plain) String() string {
	return string(t)
}

var (
	ErrOutOfRange    = errors.New("position out if range")
	ErrTitleNotFound = errors.New("title not found")
	ErrRowNotFound   = errors.New("row not found")
)

func (l Line) Width() int {
	return len(l)
}

func (l Line) Expand(width int, defaultValue Value) Line {
	line := l
	if line == nil {
		ret := make(Line, 0, width)
		for i := 0; i < width; i++ {
			ret = append(ret, defaultValue)
		}
		return ret
	}
	for wE := len(line); wE < width; wE++ {
		line = append(line, defaultValue)
	}
	return line
}

func CreateLineByStrings(lst []string) Line {
	ret := make(Line, 0, len(lst))
	for _, s := range lst {
		ret = append(ret, Plain(s))
	}
	return ret
}
