package geo

// Nearest answers which value is nearest to base.
func Nearest[T Number](base, a, b T) int {
	if a == base {
		return 0
	} else if b == base {
		return 1
	}
	dista, distb := base-a, base-b
	if dista < 0 {
		dista = -dista
	}
	if distb < 0 {
		distb = -distb
	}
	if distb < dista {
		return 1
	}
	return 0
}

func DistToPolyF(poly []PtF, pt PtF) (float64, PtF) {
	if len(poly) < 2 {
		return 0., PtF{}
	}

	prv := poly[len(poly)-1]
	foundD := 99999999.
	foundPt := PtF{}
	for _, nxt := range poly {
		seg := SegF{A: prv, B: nxt}
		if _d, _pt := DistPointToSegment(seg, pt); _d < foundD {
			foundD, foundPt = _d, _pt
		}

		prv = nxt
	}
	return foundD, foundPt
}
