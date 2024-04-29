package assign

const (
	FuzzyFloats = 1 << iota
	FuzzyInts
	FuzzyStrings
	Fuzzy = FuzzyFloats | FuzzyInts | FuzzyStrings
)
