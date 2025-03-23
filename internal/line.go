package internal

type Line struct {
	P   Vector
	Dir Vector
}

func NewLine(p, dir Vector) Line {
	return Line{p, dir.Normalize()}
}

func (l Line) Move(v Vector) Line {
	return Line{
		l.P.Add(v),
		l.Dir,
	}
}

func (l *Line) IntersectPlane(p Plane) (bool, Vector, float64) {
	den := p.NormV.Dot(l.Dir)
	if den == 0 {
		return false, Vector{}, 0
	}
	num := p.P.Dot(p.NormV) - p.NormV.Dot(l.P)
	t := num / den
	inter := l.Dir.Mul(t).Add(l.P)
	return true, inter, t
}
