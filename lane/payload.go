package lane

import (
	"errors"
)

type Payload string

type Payloads struct {
	Ver      int            `json:"v,omitempty"`
	Contents map[Tag]string `json:"c"`
}

var ErrCandidatesNotMatch = errors.New("cannot match any candidates")

func (p Payload) String() string {
	return string(p)
}

func (p Payload) Is(v string) bool {
	return p.String() == v
}

func (p Payload) In(list ...string) bool {
	n := len(list)
	if n == 0 {
		return false
	}
	for i := 0; i < n; i++ {
		if list[i] == p.String() {
			return true
		}
	}
	return false
}

func (p Payload) Select(candidates map[string]string) (string, error) {
	ret, ok := candidates[p.String()]
	if !ok {
		return "", ErrCandidatesNotMatch
	}
	return ret, nil
}
