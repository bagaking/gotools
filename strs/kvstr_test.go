package strs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVStr_ToMap(t *testing.T) {
	ret, err := KVStr("a=1,b=2.2,c=.3,d= 10e2,e =07 ,f = 0xAB2,g = ss ss, h= this is a str ,i=this is a string,j=`this is another one`,k=this-is-a-string; l='c',m='\\t',n='\\n', true").ToMap()
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
