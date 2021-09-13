package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvFile(t *testing.T) {
	err := LoadEnv("../.env")
	assert.Equal(t, nil, err, ".env file should loaded succcessfully without error")
}
