package internal

import (
	"image/color"

	"github.com/google/uuid"
)

const Eps = 1e-12

type intersectionType int

const (
	inside intersectionType = iota + 1
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
	Intersect(l *Line) *Intersection
	ID() string
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
