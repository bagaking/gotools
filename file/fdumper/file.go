package fdumper

import (
	"context"
	"io/ioutil"
	"os"
)

type (
	FileDumper struct{}
)

func (f FileDumper) Dump(ctx context.Context, filepath string, data []byte) error {
	return ioutil.WriteFile(filepath, data, os.ModePerm)
}

func (f FileDumper) Load(ctx context.Context, filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

var _ Endpoint = &FileDumper{}
