package engine

import (
	"fmt"

	"aurora.ponglehub.co.uk/pkg/webgl"
)

type Engine struct {
	gl     *webgl.WebGL
	width  int
	height int
}

func New(gl *webgl.WebGL) *Engine {
	width := gl.GetCanvasWidth()
	height := gl.GetCanvasHeight()

	return &Engine{
		gl:     gl,
		width:  width,
		height: height,
	}
}

type Drawable struct {
	Mesh         []float32
	vertexBuffer webgl.Buffer
	position     int
	Colour       []float32
	colorBuffer  webgl.Buffer
	color        int
}

type Layer struct {
	Name           string
	VertexShader   string
	FragmentShader string
	attributes     map[string]int
	program        webgl.Program
	Drawables      []*Drawable
}

type Scene struct {
	ClearColor webgl.Color
	Layers     []*Layer
}

func (e *Engine) Init(scene *Scene) error {
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

		layer.program = program

		for _, drawable := range layer.Drawables {
			drawable.vertexBuffer = gl.CreateBuffer()
			gl.BindBuffer(gl.ARRAY_BUFFER, drawable.vertexBuffer)
			gl.BufferData(gl.ARRAY_BUFFER, webgl.Float32ArrayBuffer(drawable.Mesh), gl.STATIC_DRAW)

			drawable.colorBuffer = gl.CreateBuffer()
		}
	}

	return nil
}

func (e *Engine) Render(scene *Scene) error {
	e.clear(scene)

	gl := e.gl

	for _, layer := range scene.Layers {
		gl.UseProgram(layer.program)

		for _, drawable := range layer.Drawables {
			gl.BindBuffer(gl.ARRAY_BUFFER, drawable.vertexBuffer)

			position := gl.GetAttribLocation(layer.program, "position")
			gl.VertexAttribPointer(position, 3, gl.FLOAT, false, 0, 0)
			gl.EnableVertexAttribArray(position)

			gl.BindBuffer(gl.ARRAY_BUFFER, drawable.colorBuffer)
			gl.BufferData(gl.ARRAY_BUFFER, webgl.Float32ArrayBuffer(drawable.Colour), gl.DYNAMIC_COPY)

			color := gl.GetAttribLocation(layer.program, "color")
			gl.VertexAttribPointer(color, 3, gl.FLOAT, false, 0, 0)
			gl.EnableVertexAttribArray(color)

			gl.Viewport(0, 0, e.width, e.height)
			gl.DrawArrays(gl.TRIANGLES, 0, len(drawable.Mesh)/3)
		}
	}

	return nil
}

func (e *Engine) clear(scene *Scene) {
	e.gl.ClearColor(
		scene.ClearColor.Red,
		scene.ClearColor.Green,
		scene.ClearColor.Blue,
		scene.ClearColor.Alpha,
	)
	e.gl.Clear(e.gl.COLOR_BUFFER_BIT)
}
