package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHash(t *testing.T) {
	hasher := NewBcryptHasher()
	plain := "12345"
	hashed, err := hasher.MakeHash(plain)
	assert.Empty(t, err, "error should be empty")
	assert.NotEmpty(t, hashed, "hashed string should not be empty")
}
func TestMakeHasError(t *testing.T) {
	hasher := NewBcryptHasher()
	plain := ""
	hashed, err := hasher.MakeHash(plain)
	assert.NotEmpty(t, err, "error should not be empty")
	assert.Empty(t, hashed, "hashed string should be empty")
}

func TestCompareHash(t *testing.T) {
	hasher := NewBcryptHasher()
	plain := "1231sadsad214"
	hashed, err := hasher.MakeHash(plain)
	assert.Empty(t, err, "error should be empty")
	assert.NotEmpty(t, hashed, "hashed string should not be empty")
	result, err := hasher.CompareHash(plain, *hashed)
	assert.Equal(t, true, result, "should true")
	assert.Empty(t, err, "error should be empty")
}

func TestCompareHashErr(t *testing.T) {
	hasher := NewBcryptHasher()
	plain := "1231sadsad214"
	hashed, err := hasher.MakeHash(plain)
	assert.Empty(t, err, "error should be empty")
	assert.NotEmpty(t, hashed, "hashed string should not be empty")
	rPlain := "someRandomstringplain"
	result, err := hasher.CompareHash(rPlain, *hashed)
	assert.Equal(t, false, result, "should false")
	assert.NotEmpty(t, err, "error should not be empty")
}
