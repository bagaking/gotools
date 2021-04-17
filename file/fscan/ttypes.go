package fscan

import "os"

type (
	/*
	 * FileFilter is a function that determine is a file met some conditions.
	 *
	 * the `info` param can be null when `os.Lstat` of the given path failed
	 * the `err` param will be passed to `path/filepath.Walk`
	 */
	FileFilter func(info os.FileInfo, fullPath string, searchingRoot string) (ok bool, err error)

	IScanEntry interface {
		GetRoot() string
		GetPath() string
		GetFileInfo() os.FileInfo
		GetSha1() []byte
		GetSha1Hex() (string, error)

		UpdateSha1() error
	}

	IScanner interface {
		GetRoot() string
		FireNew() (IScanResult, error)
	}

	IScanResult interface {
		IScanner

		Len() int
		UsingRelativePath() bool

		Get(path string) IScanEntry
		GetPaths() []string

		RangeFiles(fn func(pth string, se IScanEntry) error) error
		RangeDirs(fn func(pth string, se IScanEntry) error) error

		WarmUp()
		GetChildrenTable() map[string][]string
	}
)
