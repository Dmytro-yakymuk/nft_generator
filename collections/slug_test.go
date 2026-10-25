package collections_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"generator/collections"
)

func TestSlugPattern(t *testing.T) {
	valid := []string{"abc", "generator-demo", "nft-2026", "a1b", "xyz", "a23456789012345678901234567890123456789012345678901234567890123"}
	for _, slug := range valid {
		assert.True(t, collections.SlugPattern.MatchString(slug), slug)
	}

	invalid := []string{"", "ab", "-abc", "abc-", "Demo", "demo collection", "demo.site", "демо",
		"a123456789012345678901234567890123456789012345678901234567890123"}
	for _, slug := range invalid {
		assert.False(t, collections.SlugPattern.MatchString(slug), slug)
	}
}
