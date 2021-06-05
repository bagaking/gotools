package fscan

import "fmt"

func NewAndScan(searchingRoot string, usingRelativePath bool, middleWares ...FileFilter) (IScanResult, error) {
	sr := newSearchingResult(searchingRoot, usingRelativePath, middleWares)
	newSr, err := sr.FireNew()
	if err != nil {
		return nil, fmt.Errorf("scan failed, %w", err)
	}
	return newSr.(IScanResult), nil
}

func NewEmpty(searchingRoot string, usingRelativePath bool, middleWares ...FileFilter) (IScanResult, error) {
	sr := newSearchingResult(searchingRoot, usingRelativePath, middleWares)
	return sr, nil
}
