package internal

import "testing"

func BenchmarkEngine(b *testing.B) {
	engine := NewEngine()
	engine.RepositionCamera(func(f Frame) Frame {
		return f.Move(I.Mul(-2))
	})
	var cubes []Renderable
	for range 10 {
		cube, _ := NewCube(1, 1, 1)
		cubes = append(cubes, &cube)

	}

	engine.Add(cubes...)
	for b.Loop() {
		engine.Render(512, 1)

	}
}
