package engine

import (
	"fmt"

	"aurora.ponglehub.co.uk/pkg/engine/background"
	"aurora.ponglehub.co.uk/pkg/webgl"
)

type Engine struct {
	gl         *webgl.WebGL
	width      int
	height     int
	bgProgram  webgl.Program
	background []*background.Background
}

func New(gl *webgl.WebGL) *Engine {
	width := gl.GetCanvasWidth()
	height := gl.GetCanvasHeight()

	return &Engine{
		gl:         gl,
		width:      width,
		height:     height,
		background: []*background.Background{},
	}
}

func (e *Engine) Init() error {
	bgProgram, err := background.CompileProgram(e.gl)
	if err != nil {
		return fmt.Errorf("failed to compile background shaders: %+v", err)
	}

	e.bgProgram = bgProgram

	b, err := background.New("tile1")
	if err != nil {
		return fmt.Errorf("failed to load background: %+v", err)
	}

	b.Init(e.gl, e.bgProgram)
	e.background = append(e.background, b)

	return nil
}

func (e *Engine) Render() error {
	e.clear()

	e.gl.Viewport(0, 0, e.width, e.height)

	for _, bg := range e.background {
		bg.Render(e.gl, e.bgProgram)
	}

	e.gl.Flush()

	return nil
}

func (e *Engine) clear() {
	e.gl.ClearColor(
		0.95,
		0.95,
		0.95,
		1.0,
	)
	e.gl.Clear(e.gl.COLOR_BUFFER_BIT)
}
