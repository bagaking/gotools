package datatable

func (t *table) Render(titles TitleLine, showTitle bool) [][]string {
	indsMap, inds := t.getTitleMappingCache(), make([]int, 0, len(titles))

	for _, title := range titles {
		inds = append(inds, indsMap[title.Value.String()])
	}

	data := make([][]string, 0, len(t.grid)+1)
	if showTitle {
		data = append(data, titles.ToStrLst())
	}

	for i := 0; i < t.Height(); i++ {
		line := make([]string, 0, len(titles))
		for _, col := range inds {
			line = append(line, t.Get(i, col).String())
		}
		data = append(data, line)
	}

	return data
}
