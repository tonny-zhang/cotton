package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathRuleMatch(t *testing.T) {
	pr := newPathRule("/a/", nil)
	result := pr.match("/a/b")
	assert.True(t, result.IsMatch)
	assert.Nil(t, result.Params)

	assert.False(t, pr.match("/b/").IsMatch)

	pr = newPathRule("/a/:name", nil)
	result = pr.match("/a/b")
	assert.True(t, result.IsMatch)
	assert.NotNil(t, result.Params)
	assert.Equal(t, "b", result.Params["name"])
}
