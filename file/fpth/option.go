package fpth

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
func ORelativeHeader(placeholder string, val string) Option {
	return func(cfg FolderPathCfg) FolderPathCfg {
		if cfg.relativeHeader == nil {
			cfg.relativeHeader = make([]struct{ key, val string }, 0)
		}
		cfg.relativeHeader = append(cfg.relativeHeader, struct{ key, val string }{placeholder, val})
		return cfg
	}
}
