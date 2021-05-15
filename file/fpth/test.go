package fpth

import (
	"errors"
	"os"
)

var ErrIsNotDir = errors.New("the given path is not a dir")

func TestDir(pth string) error {
	stat, err := os.Stat(pth)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return ErrIsNotDir
	}
	return nil
}
