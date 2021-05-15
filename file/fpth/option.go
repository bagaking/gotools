package fpth

func OEnableHomeDir() Option {
	return func(cfg FolderPathCfg) FolderPathCfg {
		cfg.enableHomeDir = true
		return cfg
	}
}

func ORelativeGivenPath(relativeRoot string) Option {
	return func(cfg FolderPathCfg) FolderPathCfg {
		cfg.relativeRoot = relativeRoot
		return cfg
	}
}

func ORelativePWDPath() Option {
	return func(cfg FolderPathCfg) FolderPathCfg {
		cfg.relativeRoot = "."
		return cfg
	}
}
