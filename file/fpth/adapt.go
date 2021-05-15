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

	if pth[0] == '~' && cfg.enableHomeDir {
		homePth, err := GetHomePath()
		if err != nil {
			return "", err
		}
		pth = Join(homePth, pth[1:])
	}

	if pth[0] == '.' && cfg.relativeRoot != "" {
		if cfg.relativeRoot == "." {
			homePth, err := GetPWDPath()
			if err != nil {
				return "", err
			}
			pth = Join(homePth, pth[1:])
		} else {
			pth = Join(cfg.relativeRoot, pth[1:])
		}
	}

	return Clean(pth), nil
}
