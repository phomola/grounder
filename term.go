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
