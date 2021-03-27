package strs

func StartsWith(str string, prefix string) bool {
	lenPref, lenStr := len(prefix), len(str)
	if lenPref > lenStr {
		return false
	}
	return str[:lenPref] == prefix
}

func EndsWith(str string,suffix string) bool {
	n := len(str) - len(suffix)
	if n < 0 {
		return false
	}
	return str[n:] == suffix
}

func Ptr(s string) *string {
	return &s
}