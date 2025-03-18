package internal

import (
	"fmt"
	"image/color"
	"math"
)

type Triangle struct {
	P0, P1, P2 Vector
	plane      Plane
	IDGen
	color     color.Color
	edgeColor *color.Color
}

type triangleOption func(*Triangle)

func WithTriangleColor(c color.Color) triangleOption {
	return func(t *Triangle) {
		t.color = c
	}
}

func WithTriangleEdgeColor(c color.Color) triangleOption {
	return func(t *Triangle) {
		t.edgeColor = &c
	}
}

func NewTriangle(p0, p1, p2 Vector, opsf ...triangleOption) (Triangle, error) {
	// invalid if parallel diff vectors
	diff01 := p0.Sub(p1).Normalize()
	diff02 := p0.Sub(p2).Normalize()
	diff12 := p1.Sub(p2).Normalize()
	if diff01 == diff02 || diff01 == diff12 || diff02 == diff12 {
		return Triangle{}, fmt.Errorf("degenerate triangle")
	}
	// compute a plane on P0
	normV := diff01.Cross(diff02)
	t := Triangle{P0: p0, P1: p1, P2: p2, plane: NewPlane(p0, normV), IDGen: IDGen{}}
	for _, opf := range opsf {
		opf(&t)
	}
	return t, nil
}

func (t Triangle) Move(v Vector) Triangle {
	newT := t
	newT.P0 = newT.P0.Add(v)
	newT.P1 = newT.P1.Add(v)
	newT.P2 = newT.P2.Add(v)
	newT.plane = newT.plane.Move(v)
	newT.IDGen = IDGen{}
	return newT
}

func (t Triangle) Rotate(axis Line, angle Radian) Triangle {
	newT := t
	newT.P0 = newT.P0.Rotate(axis, angle)
	newT.P1 = newT.P1.Rotate(axis, angle)
	newT.P2 = newT.P2.Rotate(axis, angle)
	newT.plane = newT.plane.Rotate(axis, angle)
	newT.IDGen = IDGen{}
	return newT
}

func (t *Triangle) barycentric(p *Vector) Vector {
	// project in 2D by defining a 2D coordinate system
	xAxis := t.P0.Sub(t.P1).Normalize()
	// we want a vector orthogonal to P0 in the same plane
	// because NormV is normal to the plane, the cross will be on
	// the plane and orthoghonal to xAxis
	yAxis := xAxis.Cross(t.plane.NormV).Normalize()

	p0x, p0y := t.P0.Dot(xAxis), t.P0.Dot(yAxis)
	p1x, p1y := t.P1.Dot(xAxis), t.P1.Dot(yAxis)
	p2x, p2y := t.P2.Dot(xAxis), t.P2.Dot(yAxis)
	// we keep p in a vector for simplicity
	p = &Vector{p.Dot(xAxis), p.Dot(yAxis), 1}
	// then we need to invert a 3x3 matrix
	det := p0x*p1y + p1x*p2y + p2x*p0y - (p0x*p2y + p1x*p0y + p2x*p1y)

	u := Vector{(p1y - p2y) / det, (p2x - p1x) / det, (p1x*p2y - p2x*p1y) / det}.Dot(*p)
	v := Vector{(p2y - p0y) / det, (p0x - p2x) / det, (p2x*p0y - p0x*p2y) / det}.Dot(*p)
	w := Vector{(p0y - p1y) / det, (p1x - p0x) / det, (p0x*p1y - p1x*p0y) / det}.Dot(*p)
	return Vector{u, v, w}
}

func (t *Triangle) Intersect(l *Line) Intersection {
	hasPlaneIn, planeInter, lineT := l.IntersectPlane(t.plane)
	if !hasPlaneIn {
		return Intersection{}
	}
	barys := t.barycentric(&planeInter)
	where := inside
	for bc := range barys.Iter() {
		if math.Abs(bc) <= 3e-3 {
			if where == inside {
				where = edge
			} else {
				where = corner
			}
		} else if bc < 0 {
			where = outside
			break
		}
		// if bc < -Eps || bc > 1+Eps {
		// where = outside
		// break
		// }
	}

	color := t.color
	if t.edgeColor != nil && (where == edge || where == corner) {
		color = *t.edgeColor
	}
	// NOTE(@lberg): we can use the t from the line
	// to sign the distance so we know if a triangle in behind
	// or in front of the origin of the line
	distance := 0.
	if where != outside {
		distance = planeInter.Sub(l.P).Norm()
		if lineT < 0 {
			distance *= -1
		}
	}
	return Intersection{
		IntPoint:   planeInter,
		SignedDist: distance,
		Color:      color,
		Where:      where,
	}
}
