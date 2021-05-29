package fscan

import "fmt"

func NewAndScan(searchingRoot string, usingRelativePath bool, async bool, middleWares ...FileFilter) (IScanResult, error) {
	sr := newSearchingResult(searchingRoot, usingRelativePath, middleWares, async)
	newSr, err := sr.FireNew()
	if err != nil {
		return nil, fmt.Errorf("scan failed, %w", err)
	}
	return newSr.(IScanResult), nil
}

func NewEmpty(searchingRoot string, usingRelativePath bool, async bool, middleWares ...FileFilter) (IScanResult, error) {
	sr := newSearchingResult(searchingRoot, usingRelativePath, middleWares, async)
	return sr, nil
}
