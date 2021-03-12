package lane

type (
	Lane interface {
		inst() *lane
		Apply(tag Tag, val string) Lane
	}

	lane struct {
		Name    string         `json:"name"`
		Payload map[Tag]string `json:"setting"`
	}
)

func (l *lane) inst() *lane {
	return l
}

func New(name string) Lane {
	return &lane{
		Name:    name,
		Payload: make(map[Tag]string),
	}
}

func (l *lane) Apply(tag Tag, val string) Lane {
	tag.Write(l, val)
	return l
}
