package geo

import (
	"golang.org/x/exp/constraints"
)

type Float interface {
	constraints.Float
}

type Number interface {
	constraints.Integer | constraints.Float
}
