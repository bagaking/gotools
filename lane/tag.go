package lane

import "fmt"

type Tag string

func (tag Tag) Read(l Lane) string {
	fmt.Println(l)
	return l.inst().Payload[tag]
}

func (tag Tag) Write(l Lane, val string) {
	l.inst().Payload[tag] = val
}

func (tag Tag) Of(l Lane) Payload {
	return Payload(tag.Read(l))
}
