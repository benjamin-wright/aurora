package engine

import (
	"fmt"

	"aurora.ponglehub.co.uk/pkg/webgl"
)

type Engine struct {
	gl       *webgl.WebGL
	width    int
	height   int
	programs map[string]webgl.Program
}

func New(gl *webgl.WebGL) *Engine {
	width := gl.GetCanvasWidth()
	height := gl.GetCanvasHeight()

	return &Engine{
		gl:       gl,
		width:    width,
		height:   height,
		programs: map[string]webgl.Program{},
	}
}

type Drawable interface {
	Mesh() []float32
	Colour() []float32
}

type Layer struct {
	Name           string
	VertexShader   string
	FragmentShader string
	Drawables      []Drawable
}

type Scene struct {
	Layers []Layer
}

func (e *Engine) Init(scene *Scene) error {
	e.programs = map[string]webgl.Program{}
	gl := e.gl

	for _, layer := range scene.Layers {
		vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
		gl.ShaderSource(vertexShader, layer.VertexShader)
		gl.CompileShader(vertexShader)
		if !gl.GetShaderParameter(vertexShader, gl.COMPILE_STATUS).(bool) {
			compilationLog := gl.GetShaderInfoLog(vertexShader)
			return fmt.Errorf("compile failed (VERTEX_SHADER) %v", compilationLog)
		}

		fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
		gl.ShaderSource(fragmentShader, layer.FragmentShader)
		gl.CompileShader(fragmentShader)
		if !gl.GetShaderParameter(fragmentShader, gl.COMPILE_STATUS).(bool) {
			compilationLog := gl.GetShaderInfoLog(fragmentShader)
			return fmt.Errorf("compile failed (FRAGMENT_SHADER) %v", compilationLog)
		}

		program := gl.CreateProgram()
		gl.AttachShader(program, vertexShader)
		gl.AttachShader(program, fragmentShader)
		gl.LinkProgram(program)
		if !gl.GetProgramParameter(program, gl.LINK_STATUS).(bool) {
			return fmt.Errorf("link failed: %v", gl.GetProgramInfoLog(program))
		}

		e.programs[layer.Name] = program
	}

	return nil
}

func (e *Engine) Render(scene *Scene) error {
	e.clear()

	gl := e.gl

	for _, layer := range scene.Layers {
		gl.UseProgram(e.programs[layer.Name])

		for _, drawable := range layer.Drawables {
			program := e.programs[layer.Name]

			vertexBuffer := gl.CreateBuffer()
			gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
			gl.BufferData(gl.ARRAY_BUFFER, webgl.Float32ArrayBuffer(drawable.Mesh()), gl.STATIC_DRAW)

			position := gl.GetAttribLocation(program, "position")
			gl.VertexAttribPointer(position, 3, gl.FLOAT, false, 0, 0)
			gl.EnableVertexAttribArray(position)

			colorBuffer := gl.CreateBuffer()
			gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
			gl.BufferData(gl.ARRAY_BUFFER, webgl.Float32ArrayBuffer(drawable.Colour()), gl.STATIC_DRAW)

			color := gl.GetAttribLocation(program, "color")
			gl.VertexAttribPointer(color, 3, gl.FLOAT, false, 0, 0)
			gl.EnableVertexAttribArray(color)

			gl.Viewport(0, 0, e.width, e.height)
			gl.DrawArrays(gl.TRIANGLES, 0, len(drawable.Mesh())/3)
		}
	}

	return nil
}

func (e *Engine) clear() {
	e.gl.ClearColor(0.5, 0.5, 0.5, 0.9)
	e.gl.Clear(e.gl.COLOR_BUFFER_BIT)
}
