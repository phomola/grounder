package grounder

import (
	"strings"
)

// Term ...
type Term struct {
	Functor string
	Args    []string
}

func (t *Term) String() string {
	var sb strings.Builder
	sb.WriteString(t.Functor)
	if len(t.Args) > 0 {
		sb.WriteByte('(')
		for i, arg := range t.Args {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(arg)
		}
		sb.WriteByte(')')
	}
	return sb.String()
}

// Compare ...
func (t *Term) Compare(t2 *Term) int {
	if t.Functor < t2.Functor {
		return -1
	}
	if t.Functor > t2.Functor {
		return 1
	}
	if len(t.Args) < len(t2.Args) {
		return -1
	}
	if len(t.Args) > len(t2.Args) {
		return 1
	}
	for i, a1 := range t.Args {
		a2 := t2.Args[i]
		if c := strings.Compare(a1, a2); c != 0 {
			return c
		}
	}
	return 0
}

func (t *Term) signature() signature {
	return signature{
		functor: t.Functor,
		arity:   len(t.Args),
	}
}
