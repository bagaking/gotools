package datatable

import (
	"encoding/json"
	"errors"
	"fmt"
)

type (
	Value fmt.Stringer
	Line  []Value
	Grid  []Line
	Plain string

	Taball interface {
		GetTitle() TitleLine
		Width() int
		Height() int
		MaxRow() int

		AppendRow(line Line) error
		SetRow(row int, line Line) error
		GetRow(row int) Line

		Set(row, col int, val Value) error

		// Get returns Empty when out of range
		Get(row, col int) Value

		Query(title Title, row int) Value

		QueryByID(id Plain, resultTitle ...Title) (map[Plain]Value, error)
		StoreByID(id Plain, values map[Plain]Value) error

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
