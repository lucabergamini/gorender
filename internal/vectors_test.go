package internal

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVectorRotation(t *testing.T) {
	axis := NewLine(Zero, J)
	// NOTE(@lberg): -90 because positive rotations arround J
	// are counterclockwise with the positive J
	rot := I.Rotate(axis, DegToRad(-90.0))
	require.InDeltaSlice(t, slices.Collect(K.Iter()),
		slices.Collect(rot.Iter()), 1e-2)

	axis = NewLine(Zero, K)
	rot = I.Rotate(axis, DegToRad(90.0))
	require.InDeltaSlice(t, slices.Collect(J.Iter()),
		slices.Collect(rot.Iter()), 1e-2)
}
