package reflectool

import (
	"fmt"
	"reflect"

	"github.com/bagaking/gotools/procast"
)

type (
	Iterator func() (interface{}, error)

	ItrMapper  func(iv interface{}) (interface{}, error)
	ItrReducer func(iv interface{}, in interface{}) (interface{}, error)
)

func (itr Iterator) Next() (interface{}, error) {
	return itr()
}

func (itr Iterator) WriteTo(slicePointer interface{}, handler ...interface{}) (err error) {
	defer procast.Recover(func(e error) { err = e }, "iterator execute failed")

	vSlice, err := ToWriteableSliceValue(slicePointer)
	if err != nil {
		return fmt.Errorf("input error, %w", err)
	}

	vNewSlice := reflect.MakeSlice(vSlice.Type(), 0, vSlice.Cap())

	var mapper ItrMapper
	var reducer ItrReducer
	for _, h := range handler {
		switch t := h.(type) {
		case ItrMapper:
			mapper = t
		case ItrReducer:
			reducer = t
		}
	}

	for {
		v, err := itr.Next()
		if err != nil {
			return fmt.Errorf("itr failed, %w", err)
		}
		if v == nil {
			break
		}

		if mapper != nil {
			if v, err = mapper(v); err != nil {
				return fmt.Errorf("mapping failed, %w", err)
			}
		}

		if reducer != nil {
			if vNewSlice.Len() == 0 {
				vNewSlice = reflect.Append(vNewSlice, reflect.ValueOf(v))
			} else {
				v0 := vNewSlice.Index(0)
				if v, err = reducer(v0.Interface(), v); err != nil {
					return fmt.Errorf("reducing failed, %w", err)
				}
				v0.Set(reflect.ValueOf(v))
			}
		} else {
			vNewSlice = reflect.Append(vNewSlice, reflect.ValueOf(v))
		}

	}
	vSlice.Set(vNewSlice)
	return nil
}

func ToWriteableSliceValue(slicePointer interface{}) (*reflect.Value, error) {
	vSlicePtr := reflect.ValueOf(slicePointer)
	// to make sure it are addressable
	if vSlicePtr.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("invalid arguments, out val should be a pointer of slice %v", vSlicePtr.Type())
	}
	vSlice := vSlicePtr.Elem()
	if vSlice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid arguments, out val should be a pointer of slice %v", vSlice.Type())
	}
	return &vSlice, nil
}

func GetSliceElementType(slice interface{}) (reflect.Type, error) {
	ty := reflect.TypeOf(slice)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	if ty.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid arguments, out val should be a pointer of slice %v", ty)
	}
	return ty.Elem(), nil
}
