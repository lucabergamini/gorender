package internal

import (
	"context"
	"image"
	"math"
)

// A frame is a 3D system placed in the space
type Frame struct {
	I, J, K Vector
	P       Vector
}

var ZeroFrame = Frame{I, J, K, Zero}

func (f Frame) Move(v Vector) Frame {
	return Frame{f.I, f.J, f.K, f.P.Add(v)}
}

func (f Frame) Rotate(axis Line, angle Radian) Frame {
	return Frame{f.I.Rotate(axis, angle), f.J.Rotate(axis, angle), f.K.Rotate(axis, angle), f.P.Rotate(axis, angle)}
}

type RenderPool struct {
	poolSize int
	inChan   chan func() error
	outChan  chan error
	ctx      context.Context
}

func newRenderPool(workerNum, expectedJobs int, ctx context.Context) *RenderPool {
	return &RenderPool{
		poolSize: workerNum,
		inChan:   make(chan func() error, workerNum),
		outChan:  make(chan error, expectedJobs),
		ctx:      ctx,
	}
}

func (rp *RenderPool) Start() {
	for range rp.poolSize {
		go func() {
			for {
				select {
				case <-rp.ctx.Done():
					return
				case work := <-rp.inChan:
					rp.outChan <- work()
				}
			}
		}()
	}
}

type Camera struct {
	F    Frame
	HFov Radian
}

func (c Camera) Move(v Vector) Camera {
	newC := c
	newC.F = c.F.Move(v)
	return newC
}

func (c Camera) Rotate(axis Line, angle Radian) Camera {
	newC := c
	newC.F = c.F.Rotate(axis, angle)
	return newC
}

// RenderPerspective generates an image using ray-tracing and perspective
// perspective is achieved by using an image plane normal to camera I
// in the JK plane with sizes matching the FOV.
// and defining points there to match the pixels in the image
func (c *Camera) RenderPerspective(width int, ratio float64, objs ...Renderable) *image.RGBA {
	height := int(float64(width) / ratio)
	HFov, VFov := c.HFov, c.HFov/Radian(ratio)

	render := image.NewRGBA(image.Rect(0, 0, width, height))

	// NOTE(@lberg): this plane can be defined at any distance,
	// it does not really change things as we always cover the full section
	// of the cone (i.e. we use the FOV and not the focal distance)
	// see https://docs.blender.org/manual/en/latest/render/cameras.html

	focDis := 0.03
	start := c.F.P.Add(c.F.I.Mul(focDis))
	HOffset := focDis * math.Tan(float64(HFov)/2) / float64(width/2)
	VOffset := focDis * math.Tan(float64(VFov)/2) / float64(height/2)
	// move start to top left position
	start = start.Add(c.F.K.Mul(VOffset * float64(height) / 2)).
		Add(c.F.J.Mul(HOffset * float64(width) / 2))
	ctx, cancel := context.WithCancel(context.Background())
	pool := newRenderPool(16, width, ctx)
	pool.Start()
	defer cancel()

	for idxH := range height {
		pool.inChan <- func() error {
			for idxW := range width {
				// compute the 3D position of the pixel, we sub because of the
				// we are top left in a right system
				point := start.Sub(c.F.K.Mul(VOffset * float64(idxH))).
					Sub(c.F.J.Mul(HOffset * float64(idxW)))
				// build a line starting from camera and passing through the point
				rayLine := NewLine(c.F.P, point.Sub(c.F.P))
				var inter *Intersection
				for _, obj := range objs {
					newInter := obj.Intersect(&rayLine)
					if newInter == nil {
						continue
					}
					// if too close or behind just ignore the intersection
					if newInter.SignedDist <= focDis {
						continue
					}
					// if no current intersection or closer then current replace
					if inter == nil || newInter.SignedDist < inter.SignedDist {
						inter = newInter
					}
				}
				if inter != nil {
					render.Set(idxW, idxH, inter.Color)
				}
			}
			return nil
		}
	}

	for range width {
		<-pool.outChan
	}
	return render
}
