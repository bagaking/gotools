package fscan

import (
	"os"
	"regexp"
)

func FilterMatchRegexOfPath(strRegexMatch string) FileFilter {
	reg, err := regexp.Compile(strRegexMatch)
	return func(fi os.FileInfo, filePath string, searchingRoot string) (bool, error) {
		if err != nil {
			return false, err
		}
		if reg.MatchString(filePath) {
			return true, nil
		}
		return false, nil
	}
}

func FilterIgnoreRegexOfPath(strRegexMatch string) FileFilter {
	filterMatch := FilterMatchRegexOfPath(strRegexMatch)
	return func(fi os.FileInfo, filePath string, searchingRoot string) (bool, error) {
		matched, err := filterMatch(fi, filePath, searchingRoot)
		if err != nil {
			return false, err
		}
		return !matched, err
	}
}

func FilterIgnoreNameTableOfFile(names ...string) FileFilter {
	match := matchName(names...)
	return func(fi os.FileInfo, filePath string, searchingRoot string) (bool, error) {
		return ignoreFileOrJumpDirByName(match, filePath, searchingRoot)
	}
}

func FilterIgnorePrefixTableOfPath(prefix ...string) FileFilter {
	match := matchPrefix(prefix...)
	return func(fi os.FileInfo, filePath string, searchingRoot string) (bool, error) {
		if match(filePath) {
			return false, nil
		}
		return true, nil
	}
}

// btw. redix-tree can handle a large scale of prefix matching
func FilterIgnorePrefixTableOfFile(prefix ...string) FileFilter {
	match := matchPrefix(prefix...)
	return func(fi os.FileInfo, filePath string, searchingRoot string) (bool, error) {
		return ignoreFileOrJumpDirByName(match, filePath, searchingRoot)
	}
}

func FilterIgnoreHiddenFile() FileFilter {
	return FilterIgnorePrefixTableOfFile(".")
}
