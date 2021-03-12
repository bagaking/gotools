package lane

import "context"

type (
	Lane interface {
		inst() *lane
		Apply(tag Tag, val string) Lane
		CreateContext(ctx context.Context) context.Context
		ExtractContext(ctx context.Context) error
		TransportKey() string
	}

	lane struct {
		Name    string    `json:"name"`
		Payload *Payloads `json:"payload"`
	}
)

func (l *lane) inst() *lane {
	return l
}

func New(name string, valAndKVs ...string) Lane {
	l := &lane{
		Name: name,
		Payload: &Payloads{
			Contents: make(map[Tag]string),
		},
	}
	n := len(valAndKVs)
	if n == 1 {
		Value.WriteTo(l, valAndKVs[0])
	} else if n > 1 {
		Value.WriteTo(l, valAndKVs[0])
		valAndKVs = valAndKVs[1:]
		if n%2 == 0 {
			valAndKVs = append(valAndKVs, "")
		} else {
			n--
		}
		for i := 0; i < n; i += 2 {
			l.Payload.Contents[Tag(valAndKVs[i])] = valAndKVs[i+1]
		}
	}
	return l
}

func (l *lane) Apply(tag Tag, val string) Lane {
	tag.WriteTo(l, val)
	return l
}
