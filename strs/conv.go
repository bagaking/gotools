package strs

import (
	"strings"
)

func Conv2Snake(name string) (snake string) {
	terms := splitNameTerms(name)
	snakeTerms := make([]string, 0, len(terms))
	for _, v := range terms {
		lower := strings.ToLower(v)
		snakeTerms = append(snakeTerms, lower)
	}

	return strings.Join(snakeTerms, "_")
}

func Conv2Camel(name string) (camel string) {
	terms := splitNameTerms(name)
	camelTerms := make([]string, 0, len(terms))
	for _, v := range terms {
		lower := strings.ToLower(v)
		camelTerms = append(camelTerms, strings.ToUpper(lower[:1])+lower[1:])
	}

	return strings.Join(camelTerms, "")
}

func Conv2SnakeAndCamel(name string) (snake, camel string) {
	terms := splitNameTerms(name)
	snakeTerms := make([]string, 0, len(terms))
	camelTerms := make([]string, 0, len(terms))
	for _, v := range terms {
		lower := strings.ToLower(v)
		snakeTerms = append(snakeTerms, lower)
		camelTerms = append(camelTerms, strings.ToUpper(lower[:1])+lower[1:])
	}

	return strings.Join(snakeTerms, "_"), strings.Join(camelTerms, "")
}

func splitNameTerms(name string) []string {
	terms := make([]string, 0)
	start := -1
	termHasLower := false

	for i := 0; i < len(name); i++ {
		c := name[i]
		if !isNameTermChar(c) {
			if start >= 0 {
				terms = append(terms, name[start:i])
				start = -1
				termHasLower = false
			}
			continue
		}

		if start < 0 {
			start = i
			termHasLower = isASCIILower(c)
			continue
		}

		prev := name[i-1]
		if isASCIILower(prev) && isASCIIUpper(c) {
			terms = append(terms, name[start:i])
			start = i
			termHasLower = false
			continue
		}
		if isASCIIDigit(prev) && isASCIIUpper(c) && (termHasLower || i+1 < len(name) && isASCIILower(name[i+1])) {
			terms = append(terms, name[start:i])
			start = i
			termHasLower = false
			continue
		}
		if isASCIIUpper(prev) && isASCIIUpper(c) && i+1 < len(name) && isASCIILower(name[i+1]) {
			terms = append(terms, name[start:i])
			start = i
			termHasLower = false
			continue
		}
		if isASCIILower(c) {
			termHasLower = true
		}
	}

	if start >= 0 {
		terms = append(terms, name[start:])
	}
	return terms
}

func isNameTermChar(c byte) bool {
	return isASCIILower(c) || isASCIIUpper(c) || isASCIIDigit(c)
}

func isASCIILower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isASCIIUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isASCIIDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
