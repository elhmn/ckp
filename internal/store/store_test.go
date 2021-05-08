package store_test

import (
	"testing"

	"github.com/elhmn/ckp/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestGenerateIdempotentID(t *testing.T) {
	t.Run("returns the same hash for similar content", func(t *testing.T) {
		id1, _ := store.GenereateIdempotentID("code", "path", "comment", "alias", "solution")
		id2, _ := store.GenereateIdempotentID("code", "path", "comment", "alias", "solution")
		assert.Equal(t, id1, id2)
	})

	t.Run("returns the different hash for different content", func(t *testing.T) {
		id1, _ := store.GenereateIdempotentID("code", "path", "comment", "alias", "solution")
		id2, _ := store.GenereateIdempotentID("code1", "path2", "comment3", "alias4", "solution5")
		assert.NotEqual(t, id1, id2)
	})
}
