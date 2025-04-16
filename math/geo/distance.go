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

/*
// DistSquared answers the squared distance from the point to the segment,
// as well as the point found on the segment.
// From https://stackoverflow.com/questions/849211/shortest-distance-between-a-point-and-a-line-segment
func DistSquared(seg SegF, p PtF) (float64, PtF) {
	l2 := seg.A.DistSquared(seg.B)
	if l2 == 0 {
		return p.DistSquared(seg.A), seg.A
	}
	t := ((p.X-seg.A.X)*(seg.B.X-seg.A.X) + (p.Y-seg.A.Y)*(seg.B.Y-seg.A.Y)) / l2
	t = math.Max(0, math.Min(1, t))
	newP := PtF{X: seg.A.X + t*(seg.B.X-seg.A.X),
		Y: seg.A.Y + t*(seg.B.Y-seg.A.Y)}
	return p.DistSquared(newP), newP
}

// DistPointToSegment answers the distance from the point to the segment,
// as well as the point found on the segment.
// From https://stackoverflow.com/questions/849211/shortest-distance-between-a-point-and-a-line-segment
func DistPointToSegment(seg SegF, p PtF) (float64, PtF) {
	d, newP := DistSquared(seg, p)
	return math.Sqrt(d), newP
}
*/
