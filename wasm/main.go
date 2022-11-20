package main

import (
	"fmt"
	"log"
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

	err = eng.Init()
	if err != nil {
		return fmt.Errorf("failed to initialise game engine: %+v", err)
	}

	refreshRate := 10
	frameTime := time.Duration(int(time.Second) / refreshRate)
	for {
		last := time.Now()

		err = eng.Render()
		if err != nil {
			return fmt.Errorf("failed to render scene: %+v", err)
		}

		elapsed := time.Since(last)

		remaining := frameTime - elapsed
		if remaining > 0 {
			time.Sleep(remaining)
		}
	}

}
