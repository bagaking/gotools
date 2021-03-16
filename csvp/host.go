package csvp

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/bagaking/gotools/reflectool"
)

type Host struct {
	SlicePtrRef *reflectool.SlicePtrReflector

	index *int32
}

func NewHost(slicePtr interface{}) (*Host, error) {
	slicePtrRef, err := reflectool.NewSlicePtrReflector(slicePtr)
	if err != nil {
		return nil, fmt.Errorf("create new host failed, %w", err)
	}
	ind := int32(0)
	return &Host{
		SlicePtrRef: slicePtrRef,
		index:       &ind,
	}, nil
}

func (h *Host) GetProcess() (ind, len int) {
	i, l := atomic.LoadInt32(h.index), h.SlicePtrRef.Len()

	return int(i), l
}

func (h *Host) Rewind(i int) error {
	l := h.SlicePtrRef.Len()
	if l <= i {
		return fmt.Errorf("the index to rewind are out of range, len= %d, got %d", l, i)
	}
	atomic.StoreInt32(h.index, int32(i))
	return nil
}

func (h *Host) TakeAndForward(outItemPtr interface{}) error {
	l := h.SlicePtrRef.Len()
	if l == 0 {
		return fmt.Errorf("there are no items in the csv host")
	}
	i := int32(0)
	for {
		i = atomic.LoadInt32(h.index)
		next := (i + 1) % int32(l)
		if atomic.CompareAndSwapInt32(h.index, i, next) {
			break
		}
	}
	return h.SlicePtrRef.Read(int(i), outItemPtr)
}

func (h *Host) TakeRandom(outItemPtr interface{}) error {
	l := h.SlicePtrRef.Len()
	if l == 0 {
		return fmt.Errorf("there are no items in the csv host")
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(l)
	return h.SlicePtrRef.Read(i, outItemPtr)
}
