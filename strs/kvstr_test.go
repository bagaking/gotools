package strs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVStr_ToMap(t *testing.T) {
	kv := KVStr("a=1,b=2.2,c=.3,d= 10e2,e =07 ,f = 0xAB2,g = ss ss, h= this is a str ,i=this is a string,j=`this is another one`,k=this-is-a-string; l='c',m='\\t',n='\\n', true")
	ret, err := kv.ToMap()

	assert.Nil(t, err)
	assert.Equal(t, "1", ret["a"])
	assert.Equal(t, "2.2", ret["b"])
	assert.Equal(t, ".3", ret["c"])
	assert.Equal(t, "10e2", ret["d"])
	assert.Equal(t, "07", ret["e"])
	assert.Equal(t, "0xAB2", ret["f"])
	assert.Equal(t, "ss ss", ret["g"])
	assert.Equal(t, "this is a str", ret["h"])
	assert.Equal(t, "this is a string", ret["i"])
	assert.Equal(t, "`this is another one`", ret["j"])
	assert.Equal(t, "this-is-a-string", ret["k"])
	assert.Equal(t, "'c'", ret["l"])
	assert.Equal(t, "'\\t'", ret["m"])
	assert.Equal(t, "'\\n'", ret["n"])
	assert.Equal(t, "", ret["true"])
}

func TestKVStr_ReflectTo(t *testing.T) {
	type X struct {
		B float32
	}
	type A struct {
		X
		A int
		C string
	}
	kv := KVStr("a=1,b=2.2,c=ok")
	a := &A{}
	_, err := kv.ReflectTo(a)

	assert.Nil(t, err)
	assert.Equal(t, 1, a.A)
	assert.Equal(t, float32(2.2), a.B)
	assert.Equal(t, "ok", a.C)
}

func TestKVStr_ForEach(t *testing.T) {
	result := []struct{
		k string
		v string
	}{
		{"a", "1"},
		{"b", "2.2"},
		{"c", "ok"},
	}
	count := 0
	KVStr("a=1,b=2.2,c=ok").ForEach(func(k, v string) {
		assert.Equal(t, result[count].k, k)
		assert.Equal(t, result[count].v, v)
		count ++

	})
}