package internal

import (
	"iter"
	"math"
	"slices"
)

type Radian float64
type Degree float64

func DegToRad(d Degree) Radian {
	return Radian(d * math.Pi / 180)
}

func RadToDeg(r Radian) Degree {
	return Degree(r * 180 / math.Pi)
}

type Vector struct {
	X, Y, Z float64
}

var I = Vector{1, 0, 0}
var J = Vector{0, 1, 0}
var K = Vector{0, 0, 1}
var Zero = Vector{0, 0, 0}

func (v Vector) Add(vo Vector) Vector {
	return Vector{
		X: v.X + vo.X,
		Y: v.Y + vo.Y,
		Z: v.Z + vo.Z,
	}
}
func (v Vector) Mul(s float64) Vector {
	return Vector{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}
func (v Vector) Sub(vo Vector) Vector {
	return v.Add(vo.Neg())
}
func (v Vector) Dot(vo Vector) float64 {
	return v.X*vo.X + v.Y*vo.Y + v.Z*vo.Z
}
func (v Vector) Neg() Vector {
	return v.Mul(-1)
}
func (v Vector) Cross(vo Vector) Vector {
	return I.Mul(v.Y*vo.Z - v.Z*vo.Y).
		Add(J.Mul(v.Z*vo.X - v.X*vo.Z)).
		Add(K.Mul(v.X*vo.Y - v.Y*vo.X))
}
func (v Vector) Norm() float64 {
	val := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	return math.Sqrt(val)
}
func (v Vector) Normalize() Vector {
	norm := v.Norm()
	if norm == 0 {
		norm += Eps
	}
	return v.Mul(1 / norm)
}
func (v *Vector) Iter() iter.Seq[float64] {
	return func(yield func(float64) bool) {
		if !yield(v.X) || !yield(v.Y) || !yield(v.Z) {
			return
		}
	}
}
func (v *Vector) Slice() []float64 {
	return slices.Collect(v.Iter())
}

func (v Vector) Rotate(axis Line, angle Radian) Vector {
	// move points so that the origin passes through the axis point
	diff := axis.P.Neg()
	axis = axis.Move(diff)
	relV := v.Add(diff)
	// rodrigues
	cosA := math.Cos(float64(angle))
	sinA := math.Sin(float64(angle))
	relV = relV.Mul(cosA).
		Add(axis.Dir.Cross(relV).Mul(sinA)).
		Add(axis.Dir.Mul(axis.Dir.Dot(relV)).Mul(1 - cosA))
	// restore offset
	return relV.Sub(diff)
}

type Vector2D struct {
	X, Y float64
}
