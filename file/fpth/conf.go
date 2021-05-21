package fpth

type (
	FolderPathCfg struct {
		enableHomeDir bool

		relativeHeader []struct{ key, val string }
		relativeRoot   string
	}

	Option func(cfg FolderPathCfg) FolderPathCfg
)

var DefaultFolderPathCfg = FolderPathCfg{}

func (cfg FolderPathCfg) merge(opts []Option) FolderPathCfg {
	cp := cfg
	for _, opt := range opts {
		cp = opt(cp)
	}
	return cp
}
