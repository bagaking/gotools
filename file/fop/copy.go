package fop

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bagaking/gotools/file/fpth"
)

var (
	ErrFileConflict     = errors.New("conflicts with existing file")
	ErrFileRemoveFailed = errors.New("remove files failed")
)

func CopyDir(src string, dest string, mkDir bool, errorStop bool) error {
	if err := fpth.TestDir(src); err != nil {
		return err
	}
	if mkDir {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return err
		}
	}
	if err := fpth.TestDir(dest); err != nil {
		return err
	}
	var errs copyDirErrors
	recordErr := func(err error) error {
		if err == nil {
			return nil
		}
		if errorStop {
			return err
		}
		errs = append(errs, err)
		return nil
	}
	if err := fpth.Walk(src, func(pth string, fi os.FileInfo, err error) error {
		if err != nil {
			return recordErr(fmt.Errorf("walk %s failed: %w", pth, err))
		}
		if fi == nil {
			return recordErr(fmt.Errorf("walk %s failed: nil file info", pth))
		}
		rel, err := filepath.Rel(src, pth)
		if err != nil {
			return recordErr(fmt.Errorf("map %s relative to %s failed: %w", pth, src, err))
		}
		newPth := filepath.Join(dest, rel)
		if fi.IsDir() {
			if err = os.MkdirAll(newPth, os.ModePerm); err != nil {
				return recordErr(fmt.Errorf("mkdir %s failed: %w", newPth, err))
			}
			return nil
		}
		if err = CopyFile(pth, newPth, false); err != nil {
			return recordErr(fmt.Errorf("copy %s to %s failed: %w", pth, newPth, err))
		}
		return nil
	}); err != nil {
		if errorStop {
			return err
		}
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type copyDirErrors []error

func (errs copyDirErrors) Error() string {
	if len(errs) == 1 {
		return errs[0].Error()
	}
	var b strings.Builder
	fmt.Fprintf(&b, "copy dir failed with %d errors:", len(errs))
	for _, err := range errs {
		b.WriteString("\n - ")
		b.WriteString(err.Error())
	}
	return b.String()
}

func (errs copyDirErrors) Is(target error) bool {
	for _, err := range errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (errs copyDirErrors) Unwrap() []error {
	return []error(errs)
}

func CopyFileWithLinkRemain(src, dest string, ensureDir bool) (errRet error) { // todo: test these method with link file
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if os.ModeSymlink&si.Mode() != 0 { // symbolic link
		link, err := os.Readlink(src)
		if err != nil {
			return err
		}
		if ensureDir {
			if err = EnsureDirOfFilePth(dest); err != nil {
				return err
			}
		}
		return os.Symlink(link, dest)
	}

	return CopyFile(src, dest, ensureDir)
}

func CopyFile(src, dest string, ensureDir bool) (errRet error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if eClose := srcFile.Close(); eClose != nil {
			errRet = eClose
		}
	}()

	if ensureDir {
		if err = EnsureDirOfFilePth(dest); err != nil {
			return err
		}
	}

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	if destInfo, err := os.Stat(dest); err == nil && os.SameFile(srcInfo, destInfo) {
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	dstFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		if eClose := dstFile.Close(); eClose != nil {
			errRet = eClose
		}
	}()
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return nil
}

func SaveFile(srcStream io.Reader, dest string, override bool) (errRet error) {
	if err := os.MkdirAll(fpth.Dir(dest), os.ModePerm); err != nil {
		return fmt.Errorf("makedir failed, %w, dest= %s", err, dest)
	}

	if exist, err := fpth.PathExists(dest); err != nil {
		return fmt.Errorf("test path failed, %w, dest= %s", err, dest)
	} else if exist {
		if !override {
			return fmt.Errorf("%w, dest= %s", ErrFileConflict, dest)
		}
		if err = os.Remove(dest); err != nil {
			return fmt.Errorf("override failed, %w, dest= %s", ErrFileRemoveFailed, dest)
		}
	}

	dstFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		if eClose := dstFile.Close(); eClose != nil {
			errRet = eClose
		}
	}()
	if _, err = io.Copy(dstFile, srcStream); err != nil {
		return err
	}
	return nil
}
