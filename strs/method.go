package strs

func StartsWith(str, prefix string) bool {
	lenPref, lenStr := len(prefix), len(str)
	if lenPref > lenStr {
		return false
	}
	return str[:lenPref] == prefix
}
