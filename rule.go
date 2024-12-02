package grounder

// WeightedTermTemplate ...
type WeightedTermTemplate struct {
	Weight float64
	Term   *TermTemplate
}

// Rule ...
type Rule struct {
	ID  string
	In  []*WeightedTermTemplate
	Out []*TermTemplate
}
