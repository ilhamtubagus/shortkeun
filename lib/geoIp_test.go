package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCity_WithValidIp(t *testing.T) {
	result, err := GetCity("36.73.32.169")
	assert.Empty(t, err, "geo ip should return nil error")
	assert.NotEmpty(t, result, "geo ip should return string")
}
func TestGetCity_WithInvalidIp(t *testing.T) {
	result, err := GetCity("asdasd")
	assert.NotEmpty(t, err, "geo ip should return not nil error")
	assert.Empty(t, result, "geo ip should return nil result")
}
