package cmd

import (
	"testing"

	"github.com/elhmn/ckp/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestListScripts(t *testing.T) {
	list := fixtures.GetListWithMoreThan10Elements()

	t.Run("Returns 11 elements", func(t *testing.T) {
		got := listScripts(list, false, false, false, len(list))
		exp := fixtures.GetPrintListWithMoreThan10Elements()
		assert.Equal(t, exp, got)
	})

	t.Run("Returns 2 elements with limit of 2", func(t *testing.T) {
		got := listScripts(list, false, false, false, 2)
		exp := fixtures.GetPrintListWithLessThan2Elements()
		assert.Equal(t, exp, got)
	})

	t.Run("Returns only elements of type code", func(t *testing.T) {
		got := listScripts(list, true, false, false, len(list))
		exp := fixtures.GetPrintListOnlyCode()
		assert.Equal(t, exp, got)
	})

	t.Run("Returns only elements of type solution", func(t *testing.T) {
		got := listScripts(list, false, true, false, len(list))
		exp := fixtures.GetPrintListOnlySolution()
		assert.Equal(t, exp, got)
	})
}
