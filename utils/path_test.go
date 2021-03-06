package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanPath(t *testing.T) {
	assert.Equal(t, "/a/b", CleanPath("/a///b"))
	assert.Equal(t, "/a/b/", CleanPath("/a///b///"))
}
