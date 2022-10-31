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

func run() error {
	canvas := js.Global().Get("document").Call("getElementById", "glcanvas")

	gl, err := webgl.New(canvas)
	if err != nil {
		return fmt.Errorf("failed to create webgl context: %+v", err)
	}

	eng := engine.New(gl)

	background := engine.Drawable{
		Mesh: []float32{
			-0.5, -0.5, 0,
			0.5, -0.5, 0,
			0, 0.5, 0,
		},
		Colour: []float32{
			1, 0, 0,
			0, 1, 0,
			0, 0, 1},
	}

	scene := &engine.Scene{
		ClearColor: webgl.Color{
			Red:   0.95,
			Green: 0.95,
			Blue:  0.95,
			Alpha: 1.0,
		},
		Layers: []*engine.Layer{
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
				Drawables: []*engine.Drawable{
					&background,
				},
			},
		},
	}

	err = eng.Init(scene)
	if err != nil {
		return fmt.Errorf("failed to initialise game engine: %+v", err)
	}

	refreshRate := 60
	frameTime := time.Duration(int(time.Second) / refreshRate)
	for {
		last := time.Now()

		background.Colour[0] = rand.Float32()
		background.Colour[4] = rand.Float32()
		background.Colour[8] = rand.Float32()

		err = eng.Render(scene)
		if err != nil {
			return fmt.Errorf("failed to render scene: %+v", err)
		}

		elapsed := time.Since(last)

		log.Printf("render time: %dus", elapsed.Microseconds())

		remaining := frameTime - elapsed
		if remaining > 0 {
			time.Sleep(remaining)
		}
	}

}
