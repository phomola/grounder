package grounder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTerm(t *testing.T) {
	req := require.New(t)

	term := &Term{Functor: "a", Args: []string{"b", "c", "d"}}
	req.Equal("a(b,c,d)", term.String())
}
