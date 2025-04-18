package geo

import (
	"math"
)

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

func PointOnSegmentXY(seg Seg3dF, p PtF) (Pt3dF, float64) {
	p2, d := PointOnSegmentSquaredXY(seg, p)
	return p2, math.Sqrt(d)
}

// PointOnSegmentSquaredXY finds the nearest point from p to seg,
// but only considers the XY components, then interpolates the Z
// from the segment endpoints.
func PointOnSegmentSquaredXY(seg Seg3dF, p PtF) (Pt3dF, float64) {
	sega, segb := seg.A.XY(), seg.B.XY()
	d, _pb := DistSquared(SegF{A: sega, B: segb}, p)
	z := seg.A.Z
	if seg.A.Z != seg.B.Z {
		// Unfortunately need to sqrt it to get the right value.
		posa, posb := math.Sqrt(sega.DistSquared(_pb)), math.Sqrt(_pb.DistSquared(segb))
		pos := posa / (posa + posb)
		z = (seg.A.Z * (1. - pos)) + (seg.B.Z * pos)
	}
	return Pt3d(_pb.X, _pb.Y, z), d
}
