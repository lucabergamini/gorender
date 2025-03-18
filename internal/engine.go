package internal

import (
	"image"
	"math"
	"sync"
)

// Engine represent the 3D world and the camera observing it
type Engine struct {
	camera   Camera
	entities map[string]Renderable
	lock     sync.Mutex
}

func NewEngine() *Engine {
	return &Engine{
		camera:   Camera{ZeroFrame, math.Pi / 4},
		entities: make(map[string]Renderable),
	}
}

func (e *Engine) RepositionCamera(tr func(Frame) Frame) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.camera.F = tr(e.camera.F)
}

func (e *Engine) Add(rs ...Renderable) {
	for _, r := range rs {
		e.entities[r.ID()] = r
	}
}

func (e *Engine) Remove(r Renderable) {
	delete(e.entities, r.ID())
}

func (e *Engine) Render(width int, ratio float64) *image.RGBA {
	entities := make([]Renderable, 0, len(e.entities))
	for _, e := range e.entities {
		entities = append(entities, e)
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.camera.RenderPerspective(width, ratio, entities...)
}
