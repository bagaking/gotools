package fop

import (
	"errors"
	"fmt"
	"io"
	"os"
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
	return fpth.Walk(src, func(pth string, fi os.FileInfo, err error) error {
		if err != nil && errorStop {
			return err
		}
		newPth := strings.Replace(pth, src, dest, -1)
		if fi.IsDir() {
			if err = os.MkdirAll(newPth, os.ModePerm); err != nil && errorStop {
				return err
			}
			return nil
		}
		if err = CopyFile(pth, newPth, false); err != nil && errorStop {
			return err
		}
		return nil
	})
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

	dstFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, os.ModePerm)
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
