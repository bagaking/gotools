package fpth

import (
	"errors"
	"strings"
)

var ErrEmptyPth = errors.New("input path cannot be empty")

func Adapt(pth string, opts ...Option) (string, error) {
	cfg := DefaultFolderPathCfg.merge(opts)
	if len(pth) == 0 {
		return "", ErrEmptyPth
	}

	if pth[0] == ' ' {
		pth = strings.Trim(pth, " ")
	}

	// do placeholder first thus $PARAM can be replaced to `~` (Home) and `.` (Relative)
	// inefficient way, which based on the assumption that there are few placeholders
	for _, replacer := range cfg.replacers {
		pth = replacer(pth)
	}

	if pth[0] == '~' && cfg.enableHomeDir {
		homePth, err := GetHomePath()
		if err != nil {
			return "", err
		}
		pth = Join(homePth, pth[1:])
	}

	if cfg.relativeRoot != "" {
		result, err := joinRelative(pth, cfg.relativeRoot)
		if err != nil {
			return "", err
		}
		pth = result
	}

	return Clean(pth), nil
}

func joinRelative(pth string, relativeRoot string) (result string, err error) {
	n := len(pth)

	// invalid
	if n == 0 {
		return "", ErrEmptyPth
	}

	// count `.`
	ind := 0
	for ; ind < n && pth[ind] == '.'; ind++ {
	}

	// no need
	if ind == 0 {
		return pth, nil
	}

	result = relativeRoot
	// using pwd
	if relativeRoot == "" || relativeRoot == "." {
		result, err = GetPWDPath()
		if err != nil {
			return "", err
		}
	}

	for upfold := 1; upfold < ind; upfold++ {
		result = Join(result, "..")
	}

	return Join(result, pth[ind:]), nil
}
