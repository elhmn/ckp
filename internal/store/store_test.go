package store_test

import (
	"fmt"
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

func TestHasSensitiveDataExist(t *testing.T) {
	type Test struct {
		input   string
		expBool bool
		expWord string
	}

	t.Run("Test for `key` keyword", func(t *testing.T) {
		tests := []Test{
			{input: "je suis con", expBool: false, expWord: ""},
			{input: " je suisK ey con ", expBool: false, expWord: ""},
			{input: "je suis con key", expBool: true, expWord: "key"},
			{input: "je suis con Key", expBool: true, expWord: "key"},
			{input: "KEY je suis con", expBool: true, expWord: "key"},
			{input: " je suisKey con ", expBool: true, expWord: "key"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
				gotBool, gotWord := store.HasSensitiveData(test.input)
				assert.Equal(t, test.expBool, gotBool)
				assert.Equal(t, test.expWord, gotWord)
			})
		}
	})

	t.Run("Test for `secret` keyword", func(t *testing.T) {
		tests := []Test{
			{input: " je suisSE CRET con ", expBool: false, expWord: ""},
			{input: "je suis con secret", expBool: true, expWord: "secret"},
			{input: "je suis con SECRET", expBool: true, expWord: "secret"},
			{input: "SECRET je suis con", expBool: true, expWord: "secret"},
			{input: " je suisSECRET con ", expBool: true, expWord: "secret"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
				gotBool, gotWord := store.HasSensitiveData(test.input)
				assert.Equal(t, test.expBool, gotBool)
				assert.Equal(t, test.expWord, gotWord)
			})
		}
	})

	t.Run("Test for `auth` keyword", func(t *testing.T) {
		tests := []Test{
			{input: " je suisAU TH con ", expBool: false, expWord: ""},
			{input: "je suis con auth", expBool: true, expWord: "auth"},
			{input: "je suis con AUTH", expBool: true, expWord: "auth"},
			{input: "AUTH je suis con", expBool: true, expWord: "auth"},
			{input: " je suisAuTh con ", expBool: true, expWord: "auth"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
				gotBool, gotWord := store.HasSensitiveData(test.input)
				assert.Equal(t, test.expBool, gotBool)
				assert.Equal(t, test.expWord, gotWord)
			})
		}
	})

	t.Run("Test for `credential` keyword", func(t *testing.T) {
		tests := []Test{
			{input: " je suisCR EDENTIAL con ", expBool: false, expWord: ""},
			{input: "je suis con credential", expBool: true, expWord: "credential"},
			{input: "je suis con CREDENTIAL", expBool: true, expWord: "credential"},
			{input: "CREDENTIAL je suis con", expBool: true, expWord: "credential"},
			{input: " je suisCredEntial con ", expBool: true, expWord: "credential"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
				gotBool, gotWord := store.HasSensitiveData(test.input)
				assert.Equal(t, test.expBool, gotBool)
				assert.Equal(t, test.expWord, gotWord)
			})
		}
	})

	t.Run("Test for `creds` keyword", func(t *testing.T) {
		tests := []Test{
			{input: " je suisCR EDS con ", expBool: false, expWord: ""},
			{input: "je suis con creds", expBool: true, expWord: "creds"},
			{input: "je suis con CREDS", expBool: true, expWord: "creds"},
			{input: "CREDS je suis con", expBool: true, expWord: "creds"},
			{input: " je suisCredS con ", expBool: true, expWord: "creds"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
				gotBool, gotWord := store.HasSensitiveData(test.input)
				assert.Equal(t, test.expBool, gotBool)
				assert.Equal(t, test.expWord, gotWord)
			})
		}
	})

	t.Run("Test for `token` keyword", func(t *testing.T) {
		tests := []Test{
			{input: " je suisTO KEN con ", expBool: false, expWord: ""},
			{input: "je suis con token", expBool: true, expWord: "token"},
			{input: "je suis con TOKEN", expBool: true, expWord: "token"},
			{input: "TOKEN je suis con", expBool: true, expWord: "token"},
			{input: " je suisTokEn con ", expBool: true, expWord: "token"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
				gotBool, gotWord := store.HasSensitiveData(test.input)
				assert.Equal(t, test.expBool, gotBool)
				assert.Equal(t, test.expWord, gotWord)
			})
		}
	})

	t.Run("Test for `bearer` keyword", func(t *testing.T) {
		tests := []Test{
			{input: " je suisBEA RER con ", expBool: false, expWord: ""},
			{input: "je suis con bearer", expBool: true, expWord: "bearer"},
			{input: "je suis con BEARER", expBool: true, expWord: "bearer"},
			{input: "BEARER je suis con", expBool: true, expWord: "bearer"},
			{input: " je suisBeaRer con ", expBool: true, expWord: "bearer"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
				gotBool, gotWord := store.HasSensitiveData(test.input)
				assert.Equal(t, test.expBool, gotBool)
				assert.Equal(t, test.expWord, gotWord)
			})
		}
	})
}
