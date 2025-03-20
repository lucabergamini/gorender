package internal

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinePlaneIntersection(t *testing.T) {
	dirs := []Vector{I, J, K}
	for _, dir := range dirs {
		line := NewLine(Zero, dir)
		plane := NewPlane(dir.Mul(2), dir)
		hasIn, planIn, _ := line.IntersectPlane(plane)
		require.True(t, hasIn)
		require.Equal(t, dir.Mul(2), planIn)

		for _, dirOth := range dirs {
			if dirOth == dir {
				continue
			}
			plane := NewPlane(dirOth.Mul(2), dirOth)
			hasIn, _, _ := line.IntersectPlane(plane)
			require.False(t, hasIn)
		}
	}
}
func TestLineTriangleIntersection(t *testing.T) {
	line := NewLine(Zero, I)
	// define a base triangle on JK plane which intersects I
	startTriangle, err := NewTriangle(K.Add(I), I.Add(J).Sub(K), I.Sub(J).Sub(K))
	require.NoError(t, err)
	intersect := startTriangle.Intersect(&line)
	require.Equal(t, inside, intersect.Where)
	// move the triangle down on K to only touch I
	triangle := startTriangle.Move(K.Neg())
	intersect = triangle.Intersect(&line)
	require.Equal(t, corner, intersect.Where)
	// move the triangle up on K to intersect at the edge
	triangle = startTriangle.Move(K)
	intersect = triangle.Intersect(&line)
	require.Equal(t, edge, intersect.Where)
	// move the triangle further up, this should not intersect
	triangle = startTriangle.Move(K.Mul(2))
	intersect = triangle.Intersect(&line)
	require.Nil(t, intersect)
	// move the triangle on negative I, this should still intersect
	triangle = startTriangle.Move(I.Mul(2).Neg())
	intersect = triangle.Intersect(&line)
	require.Equal(t, inside, intersect.Where)
}

func TestFrameRotation(t *testing.T) {
	f := ZeroFrame
	f = f.Move(I).Rotate(NewLine(Zero, K), math.Pi/2)
	// we expect the point to be on J
	require.InDeltaSlice(t, J.Slice(), f.P.Slice(), 1e-2)
	// we expect the frame to have I aligned with the system J
	require.InDeltaSlice(t, J.Slice(), f.I.Slice(), 1e-2)
	// we expect the frame to have J aligned with the system -I
	negI := I.Neg()
	require.InDeltaSlice(t, negI.Slice(), f.J.Slice(), 1e-2)
	// we expect the frame to have K aligned with the system K
	require.InDeltaSlice(t, K.Slice(), f.K.Slice(), 1e-2)
}
