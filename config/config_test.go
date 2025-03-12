package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.NotEqual(t, Version, "unknown")
	assert.NotEqual(t, BuildTime, "unknown")
}
