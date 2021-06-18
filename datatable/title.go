package datatable

type (
	Title struct {
		Value Plain      `json:"value"`
		Extra TitleExtra `json:"extra"`
	}

	TitleExtra map[string]string
	TitleLine  []Title
)

func (t Title) String() string {
	return t.Value.String()
}

func (tl *TitleLine) ToStrLst() []string {
	ret := make([]string, 0, len(*tl))
	for i := range *tl {
		ret = append(ret, (*tl)[i].String())
	}
	return ret
}

func (t *Title) Match(titleVal Value, extraMatching func(extra TitleExtra) bool) bool {
	if t.Value.String() != titleVal.String() {
		return false
	}
	if extraMatching == nil {
		return true
	}
	return extraMatching(t.Extra)
}

func (t *Title) MarkID() *Title {
	if t != nil {
		t.Extra["IS_ID"] = "true"
	}
	return t
}

func (t Title) IsID() bool {
	return t.Extra["IS_ID"] == "true"
}
