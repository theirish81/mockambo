package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequiredOrRandom(t *testing.T) {
	assert.True(t, RequiredOrRandom(true))
}
