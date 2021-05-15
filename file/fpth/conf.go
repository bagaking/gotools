package fpth

type (
	FolderPathCfg struct {
		enableHomeDir bool
		relativeRoot  string
	}

	Option func(cfg FolderPathCfg) FolderPathCfg
)

func (cfg FolderPathCfg) merge(opts []Option) FolderPathCfg {
	cp := cfg
	for _, opt := range opts {
		cp = opt(cp)
	}
	return cp
}

var DefaultFolderPathCfg = FolderPathCfg{}
