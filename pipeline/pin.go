package pipeline

type Pin struct {
	Name    string
	Payload Cloner
	Policy  ClonePolicy
}

// Cloner is a copying interface.
type Cloner interface {
	Clone() Cloner
}

type ClonePolicy int

const (
	SmartClonePolicy  ClonePolicy = iota // A clone is made when one payload is going to two ore more RunInputs, othwerwise no clone.
	AlwaysClonePolicy                    // A clone is always made.
	NeverClonePolicy                     // A clone is never made.
)
