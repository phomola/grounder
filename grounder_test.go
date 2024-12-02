package grounder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddTerm(t *testing.T) {
	req := require.New(t)

	term := &WeightedTerm{
		Weight: 1.2,
		Term:   &Term{Functor: "a", Args: []string{"b", "c"}},
	}
	g := NewGrounder()

	s := g.AddTerm(term, 1)
	req.False(s)
	s = g.AddTerm(term, 1)
	req.True(s)

	req.Equal(1, g.terms.Size())
	req.Equal(1, len(g.termsByLevel))
}

func TestApplyRules(t *testing.T) {
	req := require.New(t)

	g := NewGrounder()

	g.AddRule(&Rule{
		ID: "R1",
		In: []*WeightedTermTemplate{
			{
				Weight: 0.5,
				Term:   &TermTemplate{Functor: "p", Args: []Arg{Var{Name: "x"}}},
			},
		},
		Out: []*TermTemplate{
			{Functor: "q", Args: []Arg{Var{Name: "x"}}},
		},
	})
	req.Equal(1, len(g.rules))
}
