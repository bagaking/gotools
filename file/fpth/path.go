package fpth

import (
	"os"
	"path/filepath"
)

var (
	cachedHomePath, _ = os.UserHomeDir()
	cachedPWDPath, _  = os.Getwd()

	Clean   = filepath.Clean
	Ext     = filepath.Ext
	Dir     = filepath.Dir
	Split   = filepath.Split
	Join    = filepath.Join
	Base    = filepath.Base
	Match   = filepath.Match
	Glob    = filepath.Glob
	Abs     = filepath.Abs
	IsAbs   = filepath.IsAbs
	Rel     = filepath.Rel
	Walk    = filepath.Walk
	WalkDir = filepath.WalkDir
)

func GetHomePath() (string, error) {
	if cachedHomePath != "" {
		return cachedHomePath, nil
	}
	homePth, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cachedHomePath = homePth
	return cachedHomePath, nil
}

func GetPWDPath() (string, error) {
	if cachedPWDPath != "" {
		return cachedPWDPath, nil
	}
	wdPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cachedPWDPath = wdPath
	return cachedPWDPath, nil
}
