package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}

	diff1 := &String{Value: "John Doe"}
	diff2 := &String{Value: "John Doe"}

	assert.Equal(t, hello1.HashKey(), hello2.HashKey())
	assert.Equal(t, diff1.HashKey(), diff2.HashKey())
}
