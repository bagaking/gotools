package fop

import (
	"errors"
	"fmt"
	"github.com/bagaking/gotools/procast"
	"io/fs"
	"path/filepath"
	"sync"
	//"sync/atomic"
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

var (

	DefaultWalkOption = WalkOption{
		MaxConcurrent: 1,
	}
	taskPool = sync.Pool{ New: func() interface{} { return &walkTask{} }}
)

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
	chTask := make(chan *walkTask) // give a buffer ?
	wg := sync.WaitGroup{}

	stopOrFailed := procast.NewCloseOrFailedProc(func(er error) error {
		if er == filepath.SkipDir {
			return nil
		}
		fmt.Println("=== err :", er)
		return nil // todo: if error returns, the wg may cause dead lock
	})
	stopOrFailed.GoAfterStop(func(error) {
		close(chTask)
	})
	handle := func() {
		for task, ok := <-chTask; ok; task, ok = <-chTask {
			wg.Done()

			if er := fn(task.path, task.info, task.err); er != nil {
				stopOrFailed.Fail(er)
			}

			taskPool.Put(task)
		}
	}

	for i := 0; i < int(option.MaxConcurrent); i++ {
		go handle()
	}

	errProcStop := errors.New("proc stopped")
	scan := func(path string, info fs.FileInfo, err error) error {
		if stopOrFailed.Closed() {
			return errProcStop
		}
		wg.Add(1)
		v := taskPool.Get().(*walkTask)
		v.path = path
		v.info = info
		v.err = err
		chTask <- v
		return nil
	}

	_ = procast.HoldGo(func(closer func(err error)) {
		_ = filepath.Walk(root, scan)
		closer(nil)
	})
	wg.Wait()

	return stopOrFailed.Done().Err()
}
