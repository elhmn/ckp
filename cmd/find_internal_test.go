package cmd

import (
	"fmt"
	"testing"

	"github.com/elhmn/ckp/internal/printers"
	"github.com/elhmn/ckp/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestDoesScriptContain(t *testing.T) {
	tests := []struct {
		Input string
		S     store.Script
		Exp   bool
	}{
		{Input: "solution", Exp: true, S: store.Script{
			Comment:  "my comment",
			Code:     store.Code{Content: "je suis con", Alias: "mon alias"},
			Solution: store.Solution{Content: "ma solution"},
		},
		},
		{Input: "comment", Exp: true, S: store.Script{
			Comment:  "my comment",
			Code:     store.Code{Content: "je suis con", Alias: "mon alias"},
			Solution: store.Solution{Content: "ma solution"},
		},
		},
		{Input: "je suis", Exp: true, S: store.Script{
			Comment:  "my comment",
			Code:     store.Code{Content: "je suis con", Alias: "mon alias"},
			Solution: store.Solution{Content: "ma solution"},
		},
		},
		{Input: "alias", Exp: true, S: store.Script{
			Comment:  "my comment",
			Code:     store.Code{Content: "je suis con", Alias: "mon alias"},
			Solution: store.Solution{Content: "ma solution"},
		},
		},
		{Input: "comment alias", Exp: true, S: store.Script{
			Comment:  "my comment",
			Code:     store.Code{Content: "je suis con", Alias: "mon alias"},
			Solution: store.Solution{Content: "ma solution"},
		},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d - test for \"%s\" equal %v", i, test.Input, test.Exp), func(t *testing.T) {
			assert.Equal(t, test.Exp, printers.DoesScriptContain(test.S, test.Input))
		})
	}
}
