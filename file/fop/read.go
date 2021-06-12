package fop

import (
	"fmt"
	"io/ioutil"

	"github.com/bagaking/gotools/file/fpth"
)

func ReadFile(pth string) ([]byte, error) {
	newPth, err := fpth.Adapt(pth, fpth.OEnableHomeDir())
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(newPth)
}

func MustReadFile(pth string) []byte {
	data, err := ReadFile(pth)
	if err != nil {
		panic(fmt.Errorf("read file %s failed, %w", pth, err))
	}
	return data
}
