package lane

import "fmt"

type Payload string

func (p Payload) In(list ...string) bool {
	n := len(list)
	if n == 0 {
		return false
	}
	for i := 0; i < n; i++ {
		if list[i] == string(p) {
			return true
		}
	}
	return false
}

func (p Payload) Select(candidates map[string]string) (string, error) {
	ret, ok := candidates[string(p)]
	if !ok {
		return "", fmt.Errorf("cannot find any candidates= %s", p)
	}
	return ret, nil
}
