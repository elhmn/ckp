package store_test

import (
	"testing"

	"github.com/elhmn/ckp/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestGenerateIdempotentID(t *testing.T) {
	t.Run("returns the same hash for similar content", func(t *testing.T) {
		id1, _ := store.GenereateIdempotentID("code", "comment", "alias", "solution")
		id2, _ := store.GenereateIdempotentID("code", "comment", "alias", "solution")
		assert.Equal(t, id1, id2)
	})

	t.Run("returns the different hash for different content", func(t *testing.T) {
		id1, _ := store.GenereateIdempotentID("code", "comment", "alias", "solution")
		id2, _ := store.GenereateIdempotentID("code1", "comment3", "alias4", "solution5")
		assert.NotEqual(t, id1, id2)
	})
}

func TestEntryAlreadyExist(t *testing.T) {
	t.Run("returns true when entry already exist", func(t *testing.T) {
		existingID := "my-id"
		s := store.Store{
			Scripts: []store.Script{
				store.Script{ID: existingID},
			},
		}

		assert.Equal(t, true, s.EntryAlreadyExist(existingID))
	})

	t.Run("returns false when entry does not already exist", func(t *testing.T) {
		nonExistingID := "my-new-id"
		existingID := "my-id"
		s := store.Store{
			Scripts: []store.Script{
				store.Script{ID: existingID},
			},
		}

		assert.Equal(t, false, s.EntryAlreadyExist(nonExistingID))
	})
}
