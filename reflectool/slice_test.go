package reflectool

import (
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator_WriteTo(t *testing.T) {
	i := 0
	itr := func() (interface{}, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return i, nil
	}

	ret := make([]int, 0, 10)
	err := Iterator(itr).WriteTo(ret)
	assert.NotNil(t, err, "set to slice should be failed")
	err = Iterator(itr).WriteTo(&ret)
	assert.Nil(t, err, "set to pointer of slice should be ok")
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, ret, "set to pointer of slice should be ok")
}

func TestIterator_WriteToPtr(t *testing.T) {
	i := 0

	type c struct {
		Val int
	}

	itr := func() (interface{}, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return &c{i}, nil
	}

	ret := make([]*c, 0, 10)
	err := Iterator(itr).WriteTo(ret)
	assert.NotNil(t, err, "set to slice should be failed")
	err = Iterator(itr).WriteTo(&ret)
	assert.Nil(t, err, "set to pointer of slice should be ok")
	for i := 0; i < 10; i++ {
		val := c{i + 1}
		assert.Equal(t, val, *ret[i], "value of item %d should be %v", i, val)
	}
}

func TestIterator_WriteToWithMap(t *testing.T) {
	i := 0

	type c struct {
		Val int
	}

	itr := func() (interface{}, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return i, nil
	}
	ret := make([]*c, 0, 10)
	err := Iterator(itr).WriteTo(&ret, ItrMapper(func(in interface{}) (interface{}, error) { return &c{in.(int)}, nil }))
	assert.Nil(t, err, "set to pointer of slice should be ok")
	for i := 0; i < 10; i++ {
		val := c{i + 1}
		assert.Equal(t, val, *ret[i], "value of item %d should be %v", i, val)
	}
}

func TestIterator_WriteToWithReduce(t *testing.T) {
	i := 0

	itr := func() (interface{}, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return i, nil
	}
	ret := make([]int, 0, 10)
	err := Iterator(itr).WriteTo(&ret, ItrReducer(func(a interface{}, b interface{}) (interface{}, error) { return a.(int) + b.(int), nil }))
	assert.Nil(t, err, "set to pointer of slice should be ok")
	assert.Equal(t, 55, ret[0], "value should be sum of values")
}

func TestIterator_WriteToWithMapAndReduce(t *testing.T) {
	i := 0

	type c struct {
		Val int
	}

	itr := func() (interface{}, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return &c{i}, nil
	}
	ret := make([]int, 0, 10)
	err := Iterator(itr).WriteTo(&ret,
		ItrMapper(func(in interface{}) (interface{}, error) { return in.(*c).Val, nil }),
		ItrReducer(func(a interface{}, b interface{}) (interface{}, error) { return a.(int) + b.(int), nil }),
	)
	assert.Nil(t, err, "set to pointer of slice should be ok")
	assert.Equal(t, 55, ret[0], "value should be sum of values")
}

func TestIterator_WriteToWithExitValidator(t *testing.T) {
	i := 0

	itr := func() (interface{}, error) {
		i++
		if i > 10 {
			return nil, io.EOF
		}
		return i, nil
	}
	ret := make([]int, 0, 10)
	err := Iterator(itr).WriteTo(&ret, ItrExitValidator(func(iv interface{}, err error) (bool, error) {
		if err == io.EOF {
			return true, nil
		}
		return false, err
	}))
	assert.Nil(t, err, "should exit correctly when io.EOF got")
}

func TestGetSliceElementType(t *testing.T) {
	a := make([]int, 0, 10)
	ty, err := GetSliceElementType(a)
	assert.Nil(t, err)
	assert.Equal(t, ty, reflect.TypeOf(0))
	ty, err = GetSliceElementType(&a)
	assert.Nil(t, err)
	assert.Equal(t, ty, reflect.TypeOf(0))
}
