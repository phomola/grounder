package grounder

import (
	"fmt"
	"strings"
)

// Arg ...
type Arg interface {
	fmt.Stringer
}

// String ...
type String struct {
	Value string
}

func (s String) String() string { return s.Value }

// Var ...
type Var struct {
	Name string
}

func (v Var) String() string { return "$" + v.Name }

// TermTemplate ...
type TermTemplate struct {
	Functor string
	Args    []Arg
}

func (tm *TermTemplate) String() string {
	var sb strings.Builder
	sb.WriteString(tm.Functor)
	if len(tm.Args) > 0 {
		sb.WriteByte('(')
		for i, arg := range tm.Args {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(arg.String())
		}
		sb.WriteByte(')')
	}
	return sb.String()
}

// Match ...
func (tm *TermTemplate) Match(t *Term, m map[string]string) bool {
	if tm.Functor != t.Functor || len(tm.Args) != len(t.Args) {
		return false
	}
	for i, ta := range tm.Args {
		a := t.Args[i]
		switch ta := ta.(type) {
		case Var:
			n := ta.Name
			if v, ok := m[n]; ok {
				if v != a {
					return false
				}
			} else {
				m[n] = a
			}
		case String:
			if a != ta.Value {
				return false
			}
		}
	}
	return true
}

// Ground ...
func (tm *TermTemplate) Ground(m map[string]string) (*Term, error) {
	args := make([]string, len(tm.Args))
	for i, arg := range tm.Args {
		switch arg := arg.(type) {
		case Var:
			n := arg.Name
			if v, ok := m[n]; ok {
				args[i] = v
			} else {
				return nil, fmt.Errorf("free variable '%s'", n)
			}
		case String:
			args[i] = arg.Value
		}
	}
	return &Term{Functor: tm.Functor, Args: args}, nil
}