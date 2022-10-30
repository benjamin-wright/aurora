package main

import (
	"fmt"
	"log"
	"math/rand"
	"syscall/js"
	"time"

	"aurora.ponglehub.co.uk/pkg/engine"
	"aurora.ponglehub.co.uk/pkg/webgl"
)

func main() {
	log.Println("Hello from the browser console!!")

	err := run()
	if err != nil {
		log.Printf("Error: %+v\n", err)
	}
}

type Mesh struct {
	mesh   []float32
	colour []float32
}

func (m *Mesh) Mesh() []float32 {
	return m.mesh
}

func (m *Mesh) Colour() []float32 {
	return m.colour
}

func run() error {
	canvas := js.Global().Get("document").Call("getElementById", "glcanvas")

	gl, err := webgl.New(canvas)
	if err != nil {
		return fmt.Errorf("failed to create webgl context: %+v", err)
	}

	eng := engine.New(gl)

	background := Mesh{
		mesh: []float32{
			-0.5, -0.5, 0,
			0.5, -0.5, 0,
			0, 0.5, 0,
		},
		colour: []float32{
			1, 0, 0,
			0, 1, 0,
			0, 0, 1},
	}

	scene := &engine.Scene{
		Layers: []engine.Layer{
			{
				Name: "background",
				VertexShader: `
					attribute vec3 position;
					attribute vec3 color;
					varying vec3 vColor;
					
					void main(void) {
						gl_Position = vec4(position, 1.0);
						vColor = color;
					}
				`,
				FragmentShader: `
					precision mediump float;
					varying vec3 vColor;

					void main(void) {
						gl_FragColor = vec4(vColor, 1.0);
					}
				`,
				Drawables: []engine.Drawable{
					&background,
				},
			},
		},
	}

	err = eng.Init(scene)
	if err != nil {
		return fmt.Errorf("failed to initialise game engine: %+v", err)
	}

	refreshRate := 5
	frameTime := time.Duration(int(time.Second) / refreshRate)
	last := time.Now()

	for {
		background.colour[0] = rand.Float32()
		background.colour[4] = rand.Float32()
		background.colour[8] = rand.Float32()

		err = eng.Render(scene)
		if err != nil {
			return fmt.Errorf("failed to render scene: %+v", err)
		}

		elapsed := time.Since(last)
		last = time.Now()

		remaining := frameTime - elapsed
		if remaining > 0 {
			time.Sleep(remaining)
		}
	}

}
