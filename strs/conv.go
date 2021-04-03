package strs

import (
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("([A-Z][a-z0-9]+)|([A-Z0-9]+)|([a-z0-9]+)")

func Conv2Snake(name string) (snake string) {
	terms := matchFirstCap.FindAllString(name, -1)
	snakeTerms := make([]string, 0, len(terms))
	for _, v := range terms {
		lower := strings.ToLower(v)
		snakeTerms = append(snakeTerms, lower)
	}

	return strings.Join(snakeTerms, "_")
}

func Conv2Camel(name string) (camel string) {
	terms := matchFirstCap.FindAllString(name, -1)
	camelTerms := make([]string, 0, len(terms))
	for _, v := range terms {
		lower := strings.ToLower(v)
		camelTerms = append(camelTerms, strings.ToUpper(lower[:1])+lower[1:])
	}

	return strings.Join(camelTerms, "")
}

func Conv2SnakeAndCamel(name string) (snake, camel string) {
	terms := matchFirstCap.FindAllString(name, -1)
	snakeTerms := make([]string, 0, len(terms))
	camelTerms := make([]string, 0, len(terms))
	for _, v := range terms {
		lower := strings.ToLower(v)
		snakeTerms = append(snakeTerms, lower)
		camelTerms = append(camelTerms, strings.ToUpper(lower[:1])+lower[1:])
	}

	return strings.Join(snakeTerms, "_"), strings.Join(camelTerms, "")
}
