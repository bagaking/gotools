package fpth

type (
	FolderPathCfg struct {
		replacers     []relativeFn
		enableHomeDir bool
		relativeRoot  string
	}

	Option func(cfg FolderPathCfg) FolderPathCfg

	relativeFn func(val string) string
)

var DefaultFolderPathCfg = FolderPathCfg{}

func (cfg FolderPathCfg) merge(opts []Option) FolderPathCfg {
	cp := cfg
	for _, opt := range opts {
		cp = opt(cp)
	}
	return cp
}
