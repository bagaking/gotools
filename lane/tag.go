package lane

type Tag string

const Value Tag = "_V_"

func (tag Tag) String() string {
	return string(tag)
}

func (tag Tag) ReadFrom(l Lane) string {
	return l.inst().Payload.Contents[tag]
}

func (tag Tag) WriteTo(l Lane, val string) {
	l.inst().Payload.Contents[tag] = val
}

func (tag Tag) ClearAt(l Lane) {
	delete(l.inst().Payload.Contents, tag)
}

func (tag Tag) Of(l Lane) Payload {
	return Payload(tag.ReadFrom(l))
}
