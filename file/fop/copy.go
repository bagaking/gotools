package fop

import (
	"io"
	"os"
	"strings"

	"github.com/bagaking/gotools/file/fpth"
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

	if os.ModeSymlink & si.Mode() != 0 { // symbolic link
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
		if err = os.MkdirAll(fpth.Dir(dest), os.ModePerm); err != nil {
			return err
		}
	}

	dstFile, err := os.Create(dest)
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
