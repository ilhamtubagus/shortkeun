package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHash(t *testing.T) {
	hasher := NewBcryptHash()
	plain := "12345"
	hashed, err := hasher.MakeHash(plain)
	assert.Empty(t, err, "error should be empty")
	assert.NotEmpty(t, hashed, "hashed string should not be empty")
}
func TestMakeHasError(t *testing.T) {
	hasher := NewBcryptHash()
	plain := ""
	hashed, err := hasher.MakeHash(plain)
	assert.NotEmpty(t, err, "error should not be empty")
	assert.Empty(t, hashed, "hashed string should be empty")
}

func TestCompareHash(t *testing.T) {
	hasher := NewBcryptHash()
	plain := "1231sadsad214"
	hashed, err := hasher.MakeHash(plain)
	assert.Empty(t, err, "error should be empty")
	assert.NotEmpty(t, hashed, "hashed string should not be empty")
	err = hasher.CompareHash(plain, *hashed)
	assert.Empty(t, err, "error should be empty")
}

func TestCompareHashErr(t *testing.T) {
	hasher := NewBcryptHash()
	plain := "1231sadsad214"
	hashed, err := hasher.MakeHash(plain)
	assert.Empty(t, err, "error should be empty")
	assert.NotEmpty(t, hashed, "hashed string should not be empty")
	rPlain := "someRandomstringplain"
	err = hasher.CompareHash(rPlain, *hashed)
	assert.NotEmpty(t, err, "error should not be empty")
}
