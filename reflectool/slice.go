package reflectool

import (
	"fmt"
	"reflect"

	"github.com/bagaking/gotools/procast"
)


type Iterator func() (interface{}, error)

func (itr Iterator) Next() (interface{}, error) {
	return itr()
}

func (itr Iterator) WriteTo(slice interface{}) (err error) {
	defer procast.Recover(func(e error) { err = e }, "iterator execute failed")

	vSlicePtr := reflect.ValueOf(slice)
	if vSlicePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid arguments, out val should be a pointer of slice %v", vSlicePtr.Type())
	}
	vSlice := vSlicePtr.Elem()
	if vSlice.Kind() != reflect.Slice {
		return fmt.Errorf("invalid arguments, out val should be a pointer of slice %v", vSlice.Type())
	}
	vNewSlice := reflect.MakeSlice(vSlice.Type(), 0, vSlice.Cap())

	for {
		v, err := itr.Next()
		if err != nil {
			return err
		}
		if v == nil {
			break
		}

		vNewSlice = reflect.Append(vNewSlice, reflect.ValueOf(v))
	}
	vSlice.Set(vNewSlice)
	return nil
}
