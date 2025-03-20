package main

import (
	"image"
	"lberg/gorender/internal"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

func main() {
	engine := internal.NewEngine()
	// c1, err := internal.NewQuad(internal.K,
	// 	internal.J.Add(internal.K), internal.J, internal.Zero, internal.WithQuadColor(color.Black))
	c1, err := internal.NewCube(2, 1, 3)
	if err != nil {
		panic(err)
	}

	engine.Add(&c1)
	engine.RepositionCamera(func(f internal.Frame) internal.Frame {
		return f.Move(internal.I.Mul(-5))
	})

	go func() {
		w := new(app.Window)
		w.Option(app.Title("GoRender"))
		w.Option(app.Size(512, 512))

		exitCode := 0
		if err := run(w, engine); err != nil {
			exitCode = 1
		}
		os.Exit(exitCode)
	}()
	app.Main()
}

func run(w *app.Window, engine *internal.Engine) error {
	var ops op.Ops
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	go func() {
		for {
			for range time.Tick(time.Millisecond * 30) {
				img = engine.Render(512, 1)
				w.Invalidate()
			}
		}
	}()

	imageWidget := widget.Image{Src: paint.NewImageOp(img)}
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			for {
				ev, ok := gtx.Event(
					key.Filter{Name: "W"},
					key.Filter{Name: "A"},
					key.Filter{Name: "S"},
					key.Filter{Name: "D"},
				)
				if !ok {
					break
				}
				keyEv, ok := ev.(key.Event)
				if !ok {
					break
				}
				if keyEv.State == key.Press {
					engine.RepositionCamera(func(f internal.Frame) internal.Frame {
						rot := internal.DegToRad(1)
						switch keyEv.Name {
						case "W":
							f = f.Rotate(internal.NewLine(internal.Zero, f.J), rot)
						case "A":
							f = f.Rotate(internal.NewLine(internal.Zero, internal.K), -rot)
						case "S":
							f = f.Rotate(internal.NewLine(internal.Zero, f.J), -rot)
						case "D":
							f = f.Rotate(internal.NewLine(internal.Zero, internal.K), rot)
						}
						return f
					})
				}
			}

			imageWidget.Src = paint.NewImageOp(img)
			imageWidget.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
