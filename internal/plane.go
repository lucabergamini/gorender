package internal

type Plane struct {
	P     Vector
	NormV Vector
}

func NewPlane(p, dir Vector) Plane {
	return Plane{
		P:     p,
		NormV: dir.Normalize(),
	}
}

func (p Plane) Move(v Vector) Plane {
	newP := p
	newP.P = p.P.Add(v)
	return newP
}

func (p Plane) Rotate(axis Line, angle Radian) Plane {
	newP := p
	newP.P = newP.P.Rotate(axis, angle)
	newP.NormV = newP.NormV.Rotate(axis, angle)
	return newP
}
