package fpth

import (
	"os"
	"path/filepath"
)

var (
	cachedHomePath, _ = os.UserHomeDir()
	cachedPWDPath, _  = os.Getwd()

	Clean = filepath.Clean
	Join  = filepath.Join
	Walk  = filepath.Walk
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
