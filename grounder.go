package grounder

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/fealsamh/datastructures/redblack"
)

// WeightedTerm ...
type WeightedTerm struct {
	Weight float64
	Term   *Term
}

// Compare ...
func (t *WeightedTerm) Compare(t2 *WeightedTerm) int {
	if t.Weight < t2.Weight {
		return -1
	}
	if t.Weight > t2.Weight {
		return 1
	}
	return t.Term.Compare(t2.Term)
}

type signature struct {
	functor string
	arity   int
}

type rulePos struct {
	rule *Rule
	pos  int
}

// Grounder ...
type Grounder struct {
	terms            *redblack.Set[*WeightedTerm]
	termsByLevel     map[int][]*WeightedTerm
	termsBySignature map[signature][]*WeightedTerm
	rules            []*Rule
	rulesByOutTerm   map[signature]map[rulePos]struct{}
}

// NewGrounder ...
func NewGrounder() *Grounder {
	return &Grounder{
		terms:            redblack.NewSet[*WeightedTerm](),
		termsByLevel:     make(map[int][]*WeightedTerm),
		termsBySignature: make(map[signature][]*WeightedTerm),
		rulesByOutTerm:   make(map[signature]map[rulePos]struct{}),
	}
}

// ApplyRules ...
func (g *Grounder) ApplyRules() error {
	var level int
	for {
		terms, ok := g.termsByLevel[level]
		if !ok {
			return nil
		}
		for _, term := range terms {
			if rules, ok := g.rulesByOutTerm[term.Term.signature()]; ok {
				for r := range rules {
					if err := g.applyRule(r.rule, r.pos, term, 0, make(map[string]string), nil, func(m map[string]string, out []*WeightedTerm) {
						fmt.Println("#", r.rule.ID, out)
					}); err != nil {
						return fmt.Errorf("failed to apply rule: %s %w", r.rule.ID, err)
					}
				}
			}
		}
		level++
	}
}

func (g *Grounder) applyRule(r *Rule, pos int, posTerm *WeightedTerm, i int, m map[string]string, out []*WeightedTerm, cb func(map[string]string, []*WeightedTerm)) error {
	if i == len(r.Out) {
		cb(m, out)
		return nil
	}
	tm := r.Out[i]
	if i == pos {
		vars, ok := tm.Match(posTerm.Term, m)
		if !ok {
			return errors.New("fixed pos term failed to match")
		}
		if err := g.applyRule(r, pos, posTerm, i+1, m, append(out, posTerm), cb); err != nil {
			return err
		}
		for _, n := range vars {
			delete(m, n)
		}
	} else {
		for _, t := range g.termsBySignature[tm.signature()] {
			vars, ok := tm.Match(t.Term, m)
			if ok {
				if err := g.applyRule(r, pos, posTerm, i+1, m, append(out, t), cb); err != nil {
					return err
				}
				for _, n := range vars {
					delete(m, n)
				}
			}
		}
	}
	return nil
}

func (g *Grounder) String() string {
	var sb strings.Builder
	levels := slices.Collect(maps.Keys(g.termsByLevel))
	slices.Sort(levels)
	for _, level := range levels {
		slevel := strconv.Itoa(level)
		for _, term := range g.termsByLevel[level] {
			sb.WriteString(slevel)
			sb.WriteByte(' ')
			sb.WriteString(term.Term.String())
			sb.WriteByte(' ')
			fmt.Fprintf(&sb, "%f", term.Weight)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// AddRule ...
func (g *Grounder) AddRule(r *Rule) {
	g.rules = append(g.rules, r)
	for i, t := range r.Out {
		sig := t.signature()
		m, ok := g.rulesByOutTerm[sig]
		if !ok {
			m = make(map[rulePos]struct{})
			g.rulesByOutTerm[sig] = m
		}
		m[rulePos{rule: r, pos: i}] = struct{}{}
	}
}

// AddTerm ...
func (g *Grounder) AddTerm(t *WeightedTerm, level int) bool {
	if !g.terms.Insert(t) {
		return false
	}

	terms := g.termsByLevel[level]
	g.termsByLevel[level] = append(terms, t)

	sig := t.Term.signature()
	terms = g.termsBySignature[sig]
	g.termsBySignature[sig] = append(terms, t)

	return true
}
