package fscan

import (
	"crypto/sha1"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type (
	matchHandler func(name string) bool
)

func matchPrefix(prefixes ...string) matchHandler {
	return func(name string) bool {
		for _, pref := range prefixes {
			if strings.HasPrefix(name, pref) {
				return true
			}
		}
		return false
	}
}

func matchName(names ...string) matchHandler {
	matchTable := map[string]bool{}
	for _, n := range names {
		matchTable[n] = true
	}
	return func(name string) bool {
		return matchTable[name]
	}
}

func ignoreFileOrJumpDirByName(match matchHandler, filePath string, searchingRoot string) (bool, error) {
	dir, base := path.Dir(filePath), path.Base(filePath)
	if match(base) {
		return false, nil
	}

	if dir != searchingRoot && match(dir) {
		return false, filepath.SkipDir
	}
	return true, nil
}

func calculateFileSha1(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// not caching this because of concurrency
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	return hash.Sum(nil)[:sha1.Size], nil
}
