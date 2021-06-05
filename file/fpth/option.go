package fpth

import "strings"

// OEnableHomeDir - enable the path format such as `~/a/b/c`
func OEnableHomeDir() Option {
	return func(cfg FolderPathCfg) FolderPathCfg {
		cfg.enableHomeDir = true
		return cfg
	}
}

// ORelativeGivenPath - enable the path format such as `./a/b/c`
// the `.` will be replaced with the given relativeRoot
func ORelativeGivenPath(relativeRoot string) Option {
	return func(cfg FolderPathCfg) FolderPathCfg {
		cfg.relativeRoot = relativeRoot
		return cfg
	}
}

// ORelativePWDPath - enable the path format such as `./a/b/c`
// the `.` will be replaced with the running path
func ORelativePWDPath() Option {
	return func(cfg FolderPathCfg) FolderPathCfg {
		cfg.relativeRoot = "."
		return cfg
	}
}

// ORelativeHeader - enable the path format such as `<place_holder>/a/b/c`
// the placeholder will be replaced with the running path
func ORelativeHeader(placeholder string, val string, caseIgnore bool) Option {
	return func(cfg FolderPathCfg) FolderPathCfg {
		if cfg.replacers == nil {
			cfg.replacers = make([]relativeFn, 0)
		}
		cfg.replacers = append(cfg.replacers, func(pth string) string {
			nKey, nHolder := len(pth), len(placeholder)
			if nKey < nHolder {
				return pth
			}
			match := pth[:nKey]
			if caseIgnore {
				match, placeholder = strings.ToLower(match), strings.ToLower(placeholder)
			}
			if match != placeholder {
				return pth
			}
			return val + pth[nKey:]
		})
		return cfg
	}
}
