package internal

import (
	"image/color"
	"slices"
)

type Quad struct {
	t1, t2 Triangle
	IDGen
}

type quadOption func(*Quad)

func WithQuadColor(c color.Color) quadOption {
	return func(t *Quad) {
		t.t1.color = c
		t.t2.color = c
	}
}

func WithQuadEdgeColor(c color.Color) quadOption {
	return func(t *Quad) {
		t.t1.edgeColor = &c
		t.t2.edgeColor = &c
	}
}

func NewQuad(p0, p1, p2, p3 Vector, opts ...quadOption) (Quad, error) {
	// NOTE(@lberg): 4 points can be combined in multiple 2 triangles
	// we ensure the two don't overlap by assigning the furtherst point
	// to different triangles and sharing the other two
	t1Pivot := p0
	others := []Vector{p1, p2, p3}
	bestDis, bestIdx := 0., 0
	for idx, p := range others {
		dist := t1Pivot.Sub(p).Norm()
		if dist > bestDis {
			bestDis = dist
			bestIdx = idx
		}
	}
	t2Pivot := others[bestIdx]
	others = slices.Concat(others[:bestIdx], others[bestIdx+1:])

	t1, err := NewTriangle(t1Pivot, others[0], others[1])
	if err != nil {
		return Quad{}, err
	}
	t2, err := NewTriangle(t2Pivot, others[0], others[1])
	if err != nil {
		return Quad{}, err
	}
	q := Quad{
		t1, t2, IDGen{},
	}
	for _, op := range opts {
		op(&q)
	}
	return q, nil
}

func (q Quad) Move(v Vector) Quad {
	newQ := q
	newQ.t1 = newQ.t1.Move(v)
	newQ.t2 = newQ.t2.Move(v)
	newQ.IDGen = IDGen{}
	return newQ
}

func (q Quad) Rotate(axis Line, angle Radian) Quad {
	newQ := q
	newQ.t1 = newQ.t1.Rotate(axis, angle)
	newQ.t2 = newQ.t2.Rotate(axis, angle)
	newQ.IDGen = IDGen{}
	return newQ
}

func (q *Quad) Intersect(l *Line) Intersection {
	int1 := q.t1.Intersect(l)
	int2 := q.t2.Intersect(l)
	if int1.Where == inside {
		return int1
	} else if int2.Where == inside {
		return int2
	} else if int1.Where == outside {
		return int2
	} else if int2.Where == outside {
		return int1
	}
	if int1.Where == int2.Where && int1.Where == edge {
		// use the internal color as this is an internal edge
		intIn := int1
		intIn.Color = q.t1.color
		return intIn
	}
	return int1
}

type Cube struct {
	quads []Quad
	IDGen
}

func NewCube(w, h, d float64) (Cube, error) {
	startPoints := []struct {
		start Vector
		off1  Vector
		off2  Vector
	}{
		{start: Zero.Add(K.Mul(h / 2)), off1: J.Mul(w / 2), off2: I.Mul(d / 2)},
		{start: Zero.Add(K.Neg().Mul(h / 2)), off1: J.Mul(w / 2), off2: I.Mul(d / 2)},
		{start: Zero.Add(I.Mul(d / 2)), off1: J.Mul(w / 2), off2: K.Mul(h / 2)},
		{start: Zero.Add(I.Neg().Mul(d / 2)), off1: J.Mul(w / 2), off2: K.Mul(h / 2)},
		{start: Zero.Add(J.Mul(w / 2)), off1: K.Mul(h / 2), off2: I.Mul(d / 2)},
		{start: Zero.Add(J.Neg().Mul(w / 2)), off1: K.Mul(h / 2), off2: I.Mul(d / 2)},
	}

	cube := Cube{}

	for _, sp := range startPoints {
		q, err := NewQuad(sp.start.Add(sp.off1).Add(sp.off2),
			sp.start.Add(sp.off1.Mul(-1)).Add(sp.off2),
			sp.start.Add(sp.off1).Add(sp.off2.Mul(-1)),
			sp.start.Add(sp.off1.Mul(-1)).Add(sp.off2.Mul(-1)),
			WithQuadColor(color.RGBA{255, 0, 0, 255}), WithQuadEdgeColor(color.Black))
		if err != nil {
			return Cube{}, err
		}
		cube.quads = append(cube.quads, q)

	}
	return cube, nil
}

func (c *Cube) Intersect(l *Line) Intersection {
	bestInt := Intersection{}
	for _, q := range c.quads {
		newInt := q.Intersect(l)
		if newInt.Where == outside {
			continue
		}
		if bestInt.Where == outside {
			bestInt = newInt
			continue
		}
		if newInt.SignedDist < bestInt.SignedDist {
			bestInt = newInt
		}
	}
	return bestInt
}
