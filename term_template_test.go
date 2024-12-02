package grounder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTermTemplate(t *testing.T) {
	req := require.New(t)

	tm := &TermTemplate{Functor: "a", Args: []Arg{String{"b"}, Var{"c"}, String{"d"}}}
	req.Equal("a(b,$c,d)", tm.String())
}

func TestTermMatch(t *testing.T) {
	t.Run("value clash", func(t *testing.T) {
		req := require.New(t)

		tm := &TermTemplate{Functor: "a", Args: []Arg{String{"b"}, Var{"c"}, String{"d"}}}
		term := &Term{Functor: "a", Args: []string{"b", "c", "dd"}}
		m := make(map[string]string)
		_, s := tm.Match(term, m)
		req.False(s)
	})

	t.Run("missing variable", func(t *testing.T) {
		req := require.New(t)

		tm := &TermTemplate{Functor: "a", Args: []Arg{String{"b"}, Var{"c"}, String{"d"}}}
		_, err := tm.Ground(map[string]string{"cc": "x"})
		req.NotNil(err)
		req.Equal("free variable 'c'", err.Error())
	})

	t.Run("success", func(t *testing.T) {
		req := require.New(t)

		tm := &TermTemplate{Functor: "a", Args: []Arg{String{"b"}, Var{"c"}, String{"d"}}}
		term := &Term{Functor: "a", Args: []string{"b", "c", "d"}}
		m := make(map[string]string)
		vars, s := tm.Match(term, m)
		req.True(s)
		req.Equal([]string{"c"}, vars)
		req.Equal(map[string]string{"c": "c"}, m)
		term2, err := tm.Ground(map[string]string{"c": "cc"})
		req.NoError(err)
		req.Equal("a(b,cc,d)", term2.String())
	})
}
