package datatable

type (
	Line []Value
)

func (l Line) Len() int {
	if l == nil {
		return 0
	}
	return len(l)
}

func (l Line) Get(col int) Value {
	if l == nil || col < 0 || col > l.Len() {
		return Empty
	}
	return l[col]
}

func (l Line) Gets(cols ...int) Line {
	if len(cols) == 0 {
		return l
	}
	ret := make([]Value, len(cols))
	for _, col := range cols {
		ret = append(ret, l.Get(col))
	}
	return ret
}
