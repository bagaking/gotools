package fop

import (
	"os"

	"github.com/bagaking/gotools/file/fpth"
)

var MkdirAll = os.MkdirAll

func EnsureDirOfFilePth(filePth string) error {
	return EnsureDir(fpth.Dir(filePth))
}

func EnsureDir(dirPth string) error {
	return MkdirAll(dirPth, os.ModePerm)
}
