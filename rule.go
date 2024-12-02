package grounder

// Rule ...
type Rule struct {
	ID  string
	In  []*TermTemplate
	Out []*TermTemplate
}
