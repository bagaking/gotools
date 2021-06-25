package datatable

import (
	"encoding/json"
	"fmt"
	"sync"
)

type (
	exporter struct {
		TitleLine `json:"title_line"`
		Grid      `json:"grid"`
	}

	table struct {
		titleLine         TitleLine
		grid              Grid
		titleMappingCache map[Value]int
		mu                sync.Mutex
	}
)

func (t *table) MarshalJSON() ([]byte, error) {
	return json.Marshal(exporter{
		TitleLine: t.titleLine,
		Grid:      t.grid,
	})
}

func (t *table) UnmarshalJSON(bytes []byte) error {
	exp := &exporter{}
	if err := json.Unmarshal(bytes, exp); err != nil {
		return err
	}
	t.titleLine = exp.TitleLine
	t.grid = exp.Grid
	return nil
}

func (t *table) GetTitle() TitleLine {
	return t.titleLine
}

func (t *table) Height() int {
	return len(t.grid)
}

func (t *table) MaxRow() int {
	return t.Height() - 1
}

func (t *table) Width() int {
	return len(t.titleLine)
}

func (t *table) AppendRow(line Line) error {
	return t.SetRow(t.Height(), line)
}

func (t *table) SetRow(row int, line Line) error {
	lineWidth, tableWidth := line.Width(), t.Width()
	if lineWidth > t.Width() {
		return fmt.Errorf("%w, %d > %d", ErrOutOfRange, lineWidth, tableWidth)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for max := t.MaxRow(); max < row; max++ {
		t.grid = append(t.grid, nil)
	}
	if t.grid[row] == nil {
		t.grid[row] = make(Line, 0, len(line))
	}
	lenDR := len(t.grid[row])
	for i, l := range line {
		if i < lenDR {
			t.grid[row][i] = l
		}
		t.grid[row] = append(t.grid[row], l)
	}
	return nil
}

func (t *table) GetRow(row int) Line {
	w := t.Width()
	if row >= t.Height() {
		return Line.Expand(nil, w, Empty)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.grid[row] == nil || len(t.grid[row]) < w {
		t.grid[row] = t.grid[row].Expand(w, Empty)
	}
	return t.grid[row]
}

func (t *table) Set(row, col int, val Value) error {
	if col >= t.Width() {
		return ErrOutOfRange
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	for height := t.Height(); height <= row; height++ {
		t.grid = append(t.grid, nil)
	}
	if t.grid[row] == nil || len(t.grid[row]) <= col {
		t.grid[row] = t.grid[row].Expand(col+1, Empty)
	}
	t.grid[row][col] = val
	return nil
}

func (t *table) GetLine(row int) Line {
	if t.grid == nil {
		return nil
	}
	if row < 0 || row > t.MaxRow() {
		return nil
	}
	return t.grid[row]
}

func (t *table) Gets(row int, cols ...int) Line {
	return t.GetLine(row).Gets(cols...)
}

func (t *table) Get(row, col int) Value {
	return t.GetLine(row).Get(col)
}

func (t *table) GetRowByColAndVal(col int, value Value) int {
	str := value.String()
	for i, line := range t.grid {
		if line[col].String() == str {
			return i
		}
	}
	return -1
}

func (t *table) SetByPos(row int, col int, val Value) error {
	if col >= t.Width() {
		return ErrOutOfRange
	}
	for max := t.MaxRow(); max < row; max++ {
		t.grid = append(t.grid, nil)
	}
	if t.grid[row] == nil {
		t.grid[row] = make(Line, col+1)
	}
	for len(t.grid[row]) <= col {
		t.grid[row] = append(t.grid[row], Empty)
	}
	t.grid[row][col] = val
	return nil
}

func New(title []Title) Table {
	return &table{
		titleLine: title,
		grid:      make(Grid, 0),
	}
}

func (t *table) getTitleMappingCache() map[Value]int {
	if t.titleMappingCache != nil {
		return t.titleMappingCache
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.titleMappingCache == nil {
		t.titleMappingCache = make(map[Value]int)
		for i := len(t.titleLine) - 1; i >= 0; i-- {
			t.titleMappingCache[t.titleLine[i].Value] = i
		}
	}

	return t.titleMappingCache
}
