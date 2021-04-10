package fscan

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

type (
	Entry struct {
		Root string
		Path string
		FI   os.FileInfo
		Sha1 []byte
	}
)

var _ IScanEntry = &Entry{}

func (se *Entry) String() string {
	hash, _ := se.GetSha1Hex()
	if se.FI != nil {
		if se.FI.IsDir() {
			return fmt.Sprintf("se[%s]-DIR-:%s|%d|<%s>", se.Path, se.FI.Mode(), se.FI.Size(), se.FI.ModTime().String()[:19])
		} else {
			return fmt.Sprintf("se[%s]%s:%s|%d|<%s>", se.Path, hash[:6], se.FI.Mode(), se.FI.Size(), se.FI.ModTime().String()[:19])
		}
	}
	return fmt.Sprintf("se[%s]%s:-", se.Path, hash[:6])
}

func (se *Entry) GetSha1Hex() (string, error) {
	if se.Sha1 == nil {
		if err := se.UpdateSha1(); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(se.Sha1), nil
}

func (se *Entry) UpdateSha1() error {
	pth := se.Path
	if se.Root != "" {
		pth = path.Join(se.Root, pth)
	}
	sha1bs, err := calculateFileSha1(pth)
	if err != nil {
		return err
	}
	se.Sha1 = sha1bs
	return nil
}

func (se *Entry) GetRoot() string {
	return se.Root
}

func (se *Entry) GetPath() string {
	return se.Path
}

func (se *Entry) GetFileInfo() os.FileInfo {
	return se.FI
}

func (se *Entry) GetSha1() []byte {
	return se.Sha1
}
