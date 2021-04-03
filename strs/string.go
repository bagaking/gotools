package strs

func StrOr(str string, fallback ...string) string {
	for i := 0; str == "" && i < len(fallback); i++ {
		str = fallback[i]
	}
	return str
}

func StrIfElse(ok bool, onTrue, onFalse string) string {
	if ok {
		return onTrue
	}
	return onFalse
}
