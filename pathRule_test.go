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

	pr = newPathRule("/a/:id<num>", nil)
	result = pr.match("/a/b")
	assert.False(t, result.IsMatch)
	assert.Nil(t, result.Params)
	result = pr.match("/a/123")
	assert.True(t, result.IsMatch)
	assert.NotNil(t, result.Params)
	assert.Equal(t, "123", result.Params["id"])
	result = pr.match("/a/123abc")
	assert.False(t, result.IsMatch)
	assert.Nil(t, result.Params)

	pr = newPathRule("/a/:rule{\\d+-abc}", nil)
	result = pr.match("/a/b")
	assert.False(t, result.IsMatch)
	assert.Nil(t, result.Params)
	result = pr.match("/a/123-abc")
	assert.True(t, result.IsMatch)
	assert.NotNil(t, result.Params)
	assert.Equal(t, "123-abc", result.Params["rule"])
	result = pr.match("/a/123-abctest")
	assert.False(t, result.IsMatch)
	assert.Nil(t, result.Params)

	pr = newPathRule("/a/:rule{\\d+-abc$}", nil)
	result = pr.match("/a/123-abctest")
	assert.False(t, result.IsMatch)
	assert.Nil(t, result.Params)
}
func TestIsConflictsWith(t *testing.T) {
	a := newPathRule("/a", nil)
	b := newPathRule("/:name", nil)
	assert.True(t, a.isConflictsWith(&b))

	c := newPathRule("/:id<num>", nil)
	assert.True(t, a.isConflictsWith(&c))

	d := newPathRule("/:rule{(\\d)-abc}", nil)
	assert.True(t, a.isConflictsWith(&d))

	// assert.False(t, true)
}
