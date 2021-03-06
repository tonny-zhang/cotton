package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandlerName(t *testing.T) {
	assert.Empty(t, GetHandlerName(1))
	assert.Empty(t, GetHandlerName(true))
	assert.NotEmpty(t, GetHandlerName(func() {}))
}
