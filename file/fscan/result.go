package fscan

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"

	"github.com/bagaking/gotools/debug"
)

type (
	ScanResult struct {
		// setting
		usingRelativePath bool

		// searching
		searchingRoot string
		middleWares   []FileFilter

		// result
		entries map[string]IScanEntry
		Paths   []string
	}
)

var (
	_                     IScanResult = &ScanResult{}
	_                     IScanner    = &ScanResult{}
	RecommendedConcurrent             = runtime.GOMAXPROCS(0) - 1
)

func newSearchingResult(searchingRoot string, usingRelativePath bool, middleWares []FileFilter) *ScanResult {
	sr := &ScanResult{
		searchingRoot:     searchingRoot,
		entries:           make(map[string]IScanEntry),
		middleWares:       middleWares,
		usingRelativePath: usingRelativePath,
	}
	return sr
}

func (sr *ScanResult) record(path string, fi os.FileInfo) error {
	pth, err := filepath.Rel(sr.searchingRoot, path)
	if err != nil {
		return err
	}
	root := ""
	if sr.usingRelativePath {
		path = pth
		root = sr.searchingRoot
	}

	sr.entries[path] = &Entry{
		root, path, fi, nil,
	}
	return nil
}

func (sr *ScanResult) UsingRelativePath() bool {
	return sr.usingRelativePath
}

func (sr *ScanResult) GetRoot() string {
	return sr.searchingRoot
}

func (sr *ScanResult) Get(path string) IScanEntry {
	return sr.entries[path]
}

func (sr *ScanResult) Len() int {
	return len(sr.Paths)
}

func (sr *ScanResult) GetPaths() []string {
	return sr.Paths
}

func (sr *ScanResult) RangeFiles(fn func(pth string, se IScanEntry) error) error {
	for _, pth := range sr.GetPaths() {
		se := sr.Get(pth)
		if se.GetFileInfo().IsDir() {
			continue
		}
		if err := fn(pth, se); err != nil {
			return err
		}
	}
	return nil
}

func (sr *ScanResult) RangeDirs(fn func(pth string, se IScanEntry) error) error {
	for _, pth := range sr.GetPaths() {
		se := sr.Get(pth)
		if !se.GetFileInfo().IsDir() {
			continue
		}
		if err := fn(pth, se); err != nil {
			return err
		}
	}
	return nil
}

func (sr *ScanResult) FireNew() (IScanResult, error) {
	srNew := newSearchingResult(sr.searchingRoot, sr.usingRelativePath, sr.middleWares)
	err := filepath.Walk(sr.searchingRoot, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, filter := range sr.middleWares {
			ok, err := filter(fi, path, sr.searchingRoot)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
		}
		if err := srNew.record(path, fi); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return srNew.buildPath(), nil
}

func (sr *ScanResult) buildPath() *ScanResult {
	defer debug.TimeStatisticsAndPrint("BuildPath", nil)()
	sr.Paths = make([]string, 0, len(sr.entries))
	for path := range sr.entries {
		sr.Paths = append(sr.Paths, path)
	}
	sort.Strings(sr.Paths)
	return sr
}

func (sr *ScanResult) WarmUp() {
	defer debug.TimeStatisticsAndPrint("WarmUp", nil)()
	if sr.Paths == nil || len(sr.Paths) == 0 {
		sr.buildPath()
	}
	sort.Strings(sr.Paths)
	sr.CalculateSha1Concurrently(RecommendedConcurrent, nil)
	// sr.CalculateSha1(nil)
}

func (sr *ScanResult) CalculateSha1(errorHandler func(path string, err error)) *ScanResult {
	defer debug.TimeStatisticsAndPrint("CalculateSha1", nil)()

	for pth, se := range sr.entries {
		if err := se.UpdateSha1(); err != nil && errorHandler != nil {
			errorHandler(pth, err)
		}
	}
	return sr
}

func (sr *ScanResult) CalculateSha1Concurrently(concurrent int, errorHandler func(path string, err error)) *ScanResult {
	defer debug.TimeStatisticsAndPrint("CalculateSha1Concurrently", nil)()

	taskCount := concurrent
	if taskCount < 1 {
		taskCount = 1
	}

	wg := &sync.WaitGroup{}
	total := len(sr.Paths)
	partLen := total/taskCount + 1

	fmt.Printf("[sr] start calculate sha1, fired-concurrent= %d, concurrent= %d, total= %d\n",
		RecommendedConcurrent, taskCount, total)

	for start := 0; start < total; start += partLen {
		end := start + partLen
		if end > total {
			end = total
		}
		wg.Add(1)
		go func(ss, ee int) {
			defer func() {
				fmt.Printf("[sr] calculate sha1 part finish, start= %d, end= %d", ss, ee)
				wg.Done()
			}()
			paths := sr.Paths[ss:ee]
			for _, pth := range paths {
				if err := sr.entries[pth].UpdateSha1(); err != nil && errorHandler != nil {
					errorHandler(pth, err)
				}
			}
		}(start, end)
	}
	wg.Wait()
	fmt.Printf("[sr] start calculate finished, total= %d", total) // todo: add dev log
	return sr
}
