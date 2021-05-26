package fop

import (
	"errors"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/bagaking/gotools/procast"
)

type (
	WalkOption struct {
		MaxConcurrent uint
	}

	WalkOptionFunc func(opt WalkOption) WalkOption

	walkTask struct {
		path string
		info fs.FileInfo
		err  error
	}
)

var DefaultWalkOption = WalkOption{
	MaxConcurrent: 1,
}

func WalkOptAsync(concurrent uint) WalkOptionFunc {
	return func(w WalkOption) WalkOption {
		w.MaxConcurrent = concurrent
		return w
	}
}

func (wo WalkOption) pipe(opts ...WalkOptionFunc) WalkOption {
	opt := DefaultWalkOption
	for _, optFn := range opts {
		opt = optFn(opt)
	}
	return opt
}

func Walk(root string, fn filepath.WalkFunc, opts ...WalkOptionFunc) error {
	option := DefaultWalkOption.pipe(opts...)

	if option.MaxConcurrent <= 1 {
		return filepath.Walk(root, fn)
	}

	return walkAsync(root, fn, option)
}

func walkAsync(root string, fn filepath.WalkFunc, option WalkOption) (err error) {
	wgHandle := sync.WaitGroup{}
	chTask := make(chan *walkTask) // give a buffer ?

	stopOrFailed := procast.NewCloseOrFailedProc(nil)
	stopOrFailed.GoAfterStop(func(error) {
		close(chTask)
	})

	handle := func() {
		for file, ok := <-chTask; ok; file, ok = <-chTask {
			if er := fn(file.path, file.info, file.err); er != nil {
				stopOrFailed.Fail(er)
			}
			wgHandle.Done()
		}
	}

	for i := 0; i < int(option.MaxConcurrent); i++ {
		go handle()
	}

	scan := func(path string, info fs.FileInfo, err error) error {
		if stopOrFailed.Closed() {
			return errors.New("proc stopped")
		}
		wgHandle.Add(1)
		chTask <- &walkTask{path, info, err}
		return nil
	}

	if errScan := procast.HoldGo(func(trigger func(err error)) {
		trigger(filepath.Walk(root, scan))
	}); errScan != nil {
		wgHandle.Wait()
	}

	return stopOrFailed.Done().Err()
}
