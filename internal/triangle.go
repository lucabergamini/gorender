package internal

import (
	"fmt"
	"image/color"
	"math"
)

type trianglePlaneData struct {
	plane      Plane
	axisX      Vector
	axisY      Vector
	p0, p1, p2 Vector2D
}

func newTrianglePlaneData(t *Triangle) trianglePlaneData {
	tpd := trianglePlaneData{}
	diff01 := t.P0.Sub(t.P1).Normalize()
	diff02 := t.P0.Sub(t.P2).Normalize()
	// compute a plane on P0
	normV := diff01.Cross(diff02)

	tpd.plane = NewPlane(t.P0, normV)
	tpd.axisX = t.P0.Sub(t.P1).Normalize()
	// we want a vector orthogonal to P0 in the same plane
	// because NormV is normal to the plane, the cross will be on
	// the plane and orthoghonal to xAxis
	tpd.axisY = tpd.axisX.Cross(tpd.plane.NormV).Normalize()
	tpd.p0 = Vector2D{t.P0.Dot(tpd.axisX), t.P0.Dot(tpd.axisY)}
	tpd.p1 = Vector2D{t.P1.Dot(tpd.axisX), t.P1.Dot(tpd.axisY)}
	tpd.p2 = Vector2D{t.P2.Dot(tpd.axisX), t.P2.Dot(tpd.axisY)}
	return tpd
}

type Triangle struct {
	P0, P1, P2 Vector
	IDGen
	color     color.Color
	edgeColor *color.Color
	// data required for rendering
	// which can be precomputed
	planeData trianglePlaneData
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
	t := Triangle{P0: p0, P1: p1, P2: p2, IDGen: IDGen{}}
	t.planeData = newTrianglePlaneData(&t)
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
	newT.planeData = newTrianglePlaneData(&newT)
	newT.IDGen = IDGen{}
	return newT
}

func (t Triangle) Rotate(axis Line, angle Radian) Triangle {
	newT := t
	newT.P0 = newT.P0.Rotate(axis, angle)
	newT.P1 = newT.P1.Rotate(axis, angle)
	newT.P2 = newT.P2.Rotate(axis, angle)
	newT.planeData = newTrianglePlaneData(&newT)
	newT.IDGen = IDGen{}
	return newT
}

func (t *Triangle) barycentric(p *Vector) Vector {
	tpd := t.planeData
	// we keep p in a vector for simplicity
	p = &Vector{p.Dot(tpd.axisX), p.Dot(tpd.axisY), 1}
	// then we need to invert a 3x3 matrix
	p0, p1, p2 := tpd.p0, tpd.p1, tpd.p2
	det := p0.X*p1.Y + p1.X*p2.Y + p2.X*p0.Y - (p0.X*p2.Y + p1.X*p0.Y + p2.X*p1.Y)

	u := Vector{(p1.Y - p2.Y) / det, (p2.X - p1.X) / det, (p1.X*p2.Y - p2.X*p1.Y) / det}.Dot(*p)
	v := Vector{(p2.Y - p0.Y) / det, (p0.X - p2.X) / det, (p2.X*p0.Y - p0.X*p2.Y) / det}.Dot(*p)
	w := Vector{(p0.Y - p1.Y) / det, (p1.X - p0.X) / det, (p0.X*p1.Y - p1.X*p0.Y) / det}.Dot(*p)
	return Vector{u, v, w}
}

func (t *Triangle) Intersect(l *Line) *Intersection {
	hasPlaneIn, planeInter, lineT := l.IntersectPlane(t.planeData.plane)
	if !hasPlaneIn {
		return nil
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
			return nil
		}
	}

	color := t.color
	if t.edgeColor != nil && (where == edge || where == corner) {
		color = *t.edgeColor
	}
	// NOTE(@lberg): we can use the t from the line
	// to sign the distance so we know if a triangle in behind
	// or in front of the origin of the line
	distance := planeInter.Sub(l.P).Norm()
	if lineT < 0 {
		distance *= -1
	}
	return &Intersection{
		IntPoint:   planeInter,
		SignedDist: distance,
		Color:      color,
		Where:      where,
	}
}
