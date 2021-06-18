package datatable

import (
  "encoding/json"
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
	return t.SetRow(t.MaxRow()+1, line)
}

func (t *table) SetRow(row int, line Line) error {
	if len(line) >= t.Width() {
		return ErrOutOfRange
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

func (t *table) Get(row, col int) Value {
	if t.grid == nil {
		return Empty
	}
	if row > t.MaxRow() || col >= t.Width() {
		return Empty
	}
	line := t.grid[row]
	if line == nil || col >= len(line) {
		return Empty
	}
	return line[col]
}

func (t *table) Query(title Title, row int) Value {
	col := t.titleGetOne(title, nil)
	return t.Get(row, col)
}

func (t *table) FindRow(col int, value Value) int {
	str := value.String()
	for i, line := range t.grid {
		if line[col].String() == str {
			return i
		}
	}
	return -1
}

func (t *table) QueryByID(id Plain, resultTitle ...Title) (map[Plain]Value, error) {
	row, err := t.getRowByID(id)
	if err != nil {
		return nil, err
	}

	ret := make(map[Plain]Value)
	for _, title := range resultTitle {
		col := t.titleGetOne(title.Value, nil)
		if col < 0 {
			ret[title.Value] = Empty
			continue
		}
		ret[title.Value] = t.grid[row][col]
	}

	return ret, nil
}

func (t *table) StoreByID(id Plain, values map[Plain]Value) error {
	row, err := t.getRowByID(id)
	if err != nil {
		return err
	}

	for titleValue := range values {
		col := t.titleGetOne(titleValue, nil)
		if col < 0 {
			return ErrTitleNotFound
		}
		if err = t.Set(row, col, values[titleValue]); err != nil {
			return err
		}
	}
	return nil
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

func New(title []Title) Taball {
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

func (t *table) findTitleOne(match func(title Title) bool) int {
	if match == nil {
		return -1
	}
	for i := 0; i < len(t.titleLine); i++ {
		if match(t.titleLine[i]) {
			return i
		}
	}
	return -1
}

func (t *table) titleGetOne(titleVal Value, extraMatching func(extra TitleExtra) bool) int {
	tm := t.getTitleMappingCache()
	i, ok := tm[titleVal]
	if !ok {
		return -1
	}
	if t.titleLine[i].Match(titleVal, extraMatching) {
		return i
	}
	return -1
}

func (t *table) getRowByID(id Plain) (int, error) {
	idCol := t.findTitleOne(Title.IsID)
	if idCol < 0 {
		return -1, ErrTitleNotFound
	}
	row := t.FindRow(idCol, id)
	if row < 0 {
		return -1, ErrRowNotFound
	}
	return row, nil
}
