package datatable

func (t *table) FindOnePos(q SingleQuery) (row int, col int) {
	if col = t.titleLine.SearchOneInd(func(title Title) bool {
		return title.HasTag(q.Tag)
	}); col < 0 {
		return -1, -1
	}

	if row = t.GetRowByColAndVal(col, q.Val); row < 0 {
		return -1, -1
	}
	return row, col
}

func (t *table) FindOne(q SingleQuery, selectTags ...string) Line {
	row, _ := t.FindOnePos(q)
	if len(selectTags) == 0 {
		return t.GetLine(row)
	}
	cols := t.titleLine.SearchColsByTags(selectTags...)
	if len(cols) == 0 {
		return Line{}
	}
	return t.GetLine(row).Gets(cols...)
}

func (t *table) searchTitleOne(match func(title Title) bool) int {
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

func (t *table) getColByTitle(titleVal Value, tags ...string) int {
	tm := t.getTitleMappingCache()
	i, ok := tm[titleVal]
	if !ok {
		return -1
	}
	if tags == nil || len(tags) == 0 {
		return i
	}
	for _, tag := range tags {
		if t.titleLine[i].HasTag(tag) {
			return i
		}
	}
	return -1
}
