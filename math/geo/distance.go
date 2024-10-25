package geo

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
