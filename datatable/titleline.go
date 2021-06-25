package datatable

type (
	TitleLine []Title
)

func (tl *TitleLine) ToStrLst() []string {
	ret := make([]string, 0, len(*tl))
	for i := range *tl {
		ret = append(ret, (*tl)[i].String())
	}
	return ret
}

func (tl *TitleLine) SearchOneInd(match func(title Title) bool) int {
	if match == nil {
		return -1
	}
	for i := 0; i < len(*tl); i++ {
		if match((*tl)[i]) {
			return i
		}
	}
	return -1
}

func (tl *TitleLine) SearchCols(match func(title Title) bool) (cols []int) {
	cols = []int{}
	if match == nil {
		return
	}
	for i := 0; i < len(*tl); i++ {
		if match((*tl)[i]) {
			cols = append(cols, i)
		}
	}
	return
}

func (tl *TitleLine) SearchColsByTags(tags ...string) (cols []int) {
	return tl.SearchCols(func(title Title) bool { // O(n^2)
		for _, tag := range tags {
			if title.HasTag(tag) {
				return true
			}
		}
		return false
	})
}

func (tl *TitleLine) SubTitleLine(cols []int) TitleLine {
	ret := make(TitleLine, 0, len(cols))
	for _, col := range cols {
		ret = append(ret, (*tl)[col])
	}
	return ret
}

func (tl *TitleLine) SubTitleLineByTags(tags ...string) TitleLine {
	cols := tl.SearchColsByTags(tags...)
	return tl.SubTitleLine(cols)
}
