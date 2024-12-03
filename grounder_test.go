package grounder

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mailstepcz/slice"
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
	req.True(s)
	s = g.AddTerm(term, 1)
	req.False(s)

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
				Term:   &TermTemplate{Functor: "q", Args: []Arg{Var{Name: "x"}}},
			},
		},
		Out: []*TermTemplate{
			{Functor: "r", Args: []Arg{Var{Name: "x"}}},
		},
	})
	g.AddRule(&Rule{
		ID: "R2",
		In: []*WeightedTermTemplate{
			{
				Weight: 0.25,
				Term:   &TermTemplate{Functor: "p", Args: []Arg{Var{Name: "x"}}},
			},
		},
		Out: []*TermTemplate{
			{Functor: "q", Args: []Arg{Var{Name: "x"}}},
		},
	})
	req.Equal(2, len(g.rules))

	s := g.AddTerm(&WeightedTerm{Weight: 1.5, Term: &Term{Functor: "r", Args: []string{"a"}}}, 0)
	req.True(s)

	level, ruleInstances, err := g.ApplyRules()
	req.NoError(err)
	req.Equal(3, level)

	req.Equal(2, len(ruleInstances))
	for _, ri := range ruleInstances {
		fmt.Println("#", ri.ID,
			strings.Join(slice.Fmap(func(t *WeightedTerm) string { return fmt.Sprintf("%s/%f", t.Term.String(), t.Weight) }, ri.In), " & "),
			"->",
			strings.Join(slice.Fmap(func(t *WeightedTerm) string { return fmt.Sprintf("%s/%f", t.Term.String(), t.Weight) }, ri.Out), " & "))
	}

	terms := g.String()
	req.Equal(`0 r(a) 1.500000
1 q(a) 0.750000
2 p(a) 0.187500
`, terms)
}
