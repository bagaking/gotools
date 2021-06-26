package datatable

type (
	Title struct {
		Value  Plain    `json:"value"`
		tags   []string `json:"tags"`
		tagMap map[string]bool
	}
)

func NewTitle(value string, tags ...string) *Title {
	ret := &Title{
		Value:  Plain(value),
		tags:   tags,
		tagMap: make(map[string]bool),
	}
	for _, tag := range tags {
		ret.tagMap[tag] = true
	}
	return ret
}

func (t *Title) DerivativeCovered(value string) *Title {
	return NewTitle(value, t.tags...)
}

func (t *Title) DerivativeAddOn(addOn string) *Title {
	return NewTitle(t.Value.String()+addOn, t.tags...)
}

func (t *Title) String() string {
	if t == nil {
		return ""
	}
	return t.Value.String()
}

func (t *Title) Mark(tag string) *Title {
	if t == nil {
		return t
	}
	if t.tags == nil {
		t.tags = make([]string, 0)
	}
	if t.tagMap == nil {
		t.tagMap = make(map[string]bool)
	}
	t.tags = append(t.tags, tag)
	t.tagMap[tag] = true
	return t
}

func (t *Title) HasTag(tag string) bool {
	if t == nil {
		return false
	}
	if t.tagMap == nil {
		return false
	}
	return t.tagMap[tag]
}
