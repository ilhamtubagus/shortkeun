package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomGenerator(t *testing.T) {
	length := 5
	rand := RandString(length)
	assert.NotEmpty(t, rand, "should not be empty")
	assert.Equal(t, length, len(rand), "length should be match")
}
