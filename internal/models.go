package internal

import (
	"image/color"

	"github.com/google/uuid"
)

const Eps = 1e-12

type intersectionType int

const (
	outside intersectionType = iota
	inside
	corner
	edge
)

type Intersection struct {
	IntPoint   Vector
	SignedDist float64
	Color      color.Color
	Where      intersectionType
}

type Renderable interface {
	Intersect(l *Line) Intersection
	ID() string
}

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

type IDGen struct {
	id *string
}

func (ig *IDGen) ID() string {
	if ig.id != nil {
		return *ig.id
	}
	ig.id = new(string)
	*ig.id = uuid.NewString()
	return *ig.id
}
