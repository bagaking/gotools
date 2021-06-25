package datatable

type (
	Title struct {
		Value Value `json:"value"`
		Tags  Tags  `json:"tags"`
	}

	Tags map[string]bool
)

func (t Title) String() string {
	return t.Value.String()
}

func (t *Title) Distinct(match func(extra Tags) bool) bool {
	if match == nil {
		return false
	}
	return match(t.Tags)
}

func (t *Title) Match(titleVal Value, extraMatching func(extra Tags) bool) bool {
	if t.Value.String() != titleVal.String() {
		return false
	}
	if extraMatching == nil {
		return true
	}
	return extraMatching(t.Tags)
}

func (t *Title) MarkTag(tag string) *Title {
	if t != nil && t.Tags != nil {
		t.Tags[tag] = true
	}
	return t
}

func (t Title) HasTag(tag string) bool {
	return t.Tags[tag]
}
