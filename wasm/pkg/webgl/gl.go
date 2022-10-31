package webgl

import (
	"errors"
	"syscall/js"
)

type BufferType int
type BufferUsage int
type BufferMask int
type Capacity int
type DrawMode int
type ShaderType int
type Type int
type ProgramParameter int
type ShaderParameter int

type Buffer js.Value
type Shader js.Value
type Location js.Value
type Program js.Value

var (
	float32Array = js.Global().Get("Float32Array")
	uint8Array   = js.Global().Get("Uint8Array")
)

type WebGL struct {
	canvas js.Value
	gl     js.Value

	ARRAY_BUFFER                                                                  BufferType
	STATIC_DRAW, DYNAMIC_COPY, STREAM_READ                                        BufferUsage
	COLOR_BUFFER_BIT                                                              BufferMask
	VERTEX_SHADER, FRAGMENT_SHADER                                                ShaderType
	COMPILE_STATUS                                                                ShaderParameter
	LINK_STATUS, VALIDATE_STATUS                                                  ProgramParameter
	POINTS, LINE_STRIP, LINE_LOOP, LINES, TRIANGLE_STRIP, TRIANGLE_FAN, TRIANGLES DrawMode
	FLOAT, UNSIGNED_BYTE, UNSIGNED_SHORT, UNSIGNED_INT                            Type
}

func New(canvas js.Value) (*WebGL, error) {
	gl := canvas.Call("getContext", "webgl2")
	if gl.IsNull() {
		return nil, errors.New("WebGL is not supported")
	}

	return &WebGL{
		canvas:       canvas,
		gl:           gl,
		ARRAY_BUFFER: BufferType(gl.Get("ARRAY_BUFFER").Int()),

		STATIC_DRAW:  BufferUsage(gl.Get("STATIC_DRAW").Int()),
		DYNAMIC_COPY: BufferUsage(gl.Get("DYNAMIC_COPY").Int()),
		STREAM_READ:  BufferUsage(gl.Get("STREAM_READ").Int()),

		COLOR_BUFFER_BIT: BufferMask(gl.Get("COLOR_BUFFER_BIT").Int()),
		VERTEX_SHADER:    ShaderType(gl.Get("VERTEX_SHADER").Int()),
		FRAGMENT_SHADER:  ShaderType(gl.Get("FRAGMENT_SHADER").Int()),

		COMPILE_STATUS:  ShaderParameter(gl.Get("COMPILE_STATUS").Int()),
		LINK_STATUS:     ProgramParameter(gl.Get("LINK_STATUS").Int()),
		VALIDATE_STATUS: ProgramParameter(gl.Get("VALIDATE_STATUS").Int()),

		FLOAT:          Type(gl.Get("FLOAT").Int()),
		UNSIGNED_BYTE:  Type(gl.Get("UNSIGNED_BYTE").Int()),
		UNSIGNED_SHORT: Type(gl.Get("UNSIGNED_SHORT").Int()),
		UNSIGNED_INT:   Type(gl.Get("UNSIGNED_INT").Int()),

		POINTS:         DrawMode(gl.Get("POINTS").Int()),
		LINE_STRIP:     DrawMode(gl.Get("LINE_STRIP").Int()),
		LINE_LOOP:      DrawMode(gl.Get("LINE_LOOP").Int()),
		LINES:          DrawMode(gl.Get("LINES").Int()),
		TRIANGLE_STRIP: DrawMode(gl.Get("TRIANGLE_STRIP").Int()),
		TRIANGLE_FAN:   DrawMode(gl.Get("TRIANGLE_FAN").Int()),
		TRIANGLES:      DrawMode(gl.Get("TRIANGLES").Int()),
	}, nil
}

func (gl *WebGL) GetCanvasWidth() int {
	return js.Value(gl.canvas).Get("clientWidth").Int()
}

func (gl *WebGL) GetCanvasHeight() int {
	return js.Value(gl.canvas).Get("clientHeight").Int()
}

func (gl *WebGL) Viewport(x1, y1, x2, y2 int) {
	gl.gl.Call("viewport", x1, y1, x2, y2)
}

func (gl *WebGL) CreateBuffer() Buffer {
	return Buffer(gl.gl.Call("createBuffer"))
}

func (gl *WebGL) DeleteBuffer(buf Buffer) {
	gl.gl.Call("deleteBuffer", js.Value(buf))
}

func (gl *WebGL) BindBuffer(t BufferType, buf Buffer) {
	gl.gl.Call("bindBuffer", int(t), js.Value(buf))
}

func (gl *WebGL) BufferData(t BufferType, data BufferData, usage BufferUsage) {
	bin := data.Bytes()
	dataJS := uint8Array.New(len(bin))
	js.CopyBytesToJS(dataJS, bin)
	gl.gl.Call("bufferData", int(t), dataJS, int(usage))
}

func (gl *WebGL) Clear(mask BufferMask) {
	gl.gl.Call("clear", int(mask))
}

func (gl *WebGL) ClearColor(r, g, b, a float32) {
	gl.gl.Call("clearColor", r, g, b, a)
}

func (gl *WebGL) ClearDepth(d float32) {
	gl.gl.Call("clearDepth", d)
}

func (gl *WebGL) Enable(c Capacity) {
	gl.gl.Call("enable", int(c))
}

func (gl *WebGL) CreateShader(t ShaderType) Shader {
	return Shader(gl.gl.Call("createShader", int(t)))
}

func (gl *WebGL) ShaderSource(s Shader, src string) {
	gl.gl.Call("shaderSource", js.Value(s), src)
}

func (gl *WebGL) CompileShader(s Shader) {
	gl.gl.Call("compileShader", js.Value(s))
}

func (gl *WebGL) GetShaderParameter(s Shader, param ShaderParameter) interface{} {
	v := gl.gl.Call("getShaderParameter", js.Value(s), int(param))
	switch param {
	case gl.COMPILE_STATUS:
		return v.Bool()
	}
	return nil
}

func (gl *WebGL) GetShaderInfoLog(s Shader) string {
	return gl.gl.Call("getShaderInfoLog", js.Value(s)).String()
}

func (gl *WebGL) GetAttribLocation(p Program, name string) int {
	return gl.gl.Call("getAttribLocation", js.Value(p), name).Int()
}

func (gl *WebGL) VertexAttribPointer(i, size int, typ Type, normalized bool, stride, offset int) {
	gl.gl.Call("vertexAttribPointer", i, size, int(typ), normalized, stride, offset)
}

func (gl *WebGL) EnableVertexAttribArray(i int) {
	gl.gl.Call("enableVertexAttribArray", i)
}

func (gl *WebGL) DisableVertexAttribArray(i int) {
	gl.gl.Call("disableVertexAttribArray", i)
}

func (gl *WebGL) GetUniformLocation(p Program, name string) Location {
	return Location(gl.gl.Call("getUniformLocation", js.Value(p), name))
}

func (gl *WebGL) CreateProgram() Program {
	return Program(gl.gl.Call("createProgram"))
}

func (gl *WebGL) GetProgramParameter(p Program, param ProgramParameter) interface{} {
	v := gl.gl.Call("getProgramParameter", js.Value(p), int(param))
	switch param {
	case gl.LINK_STATUS, gl.VALIDATE_STATUS:
		return v.Bool()
	}
	return nil
}

func (gl *WebGL) GetProgramInfoLog(p Program) string {
	return gl.gl.Call("getProgramInfoLog", js.Value(p)).String()
}

func (gl *WebGL) AttachShader(p Program, s Shader) {
	gl.gl.Call("attachShader", js.Value(p), js.Value(s))
}

func (gl *WebGL) LinkProgram(p Program) {
	gl.gl.Call("linkProgram", js.Value(p))
}

func (gl *WebGL) UseProgram(p Program) {
	gl.gl.Call("useProgram", js.Value(p))
}

func (gl *WebGL) DrawArrays(mode DrawMode, i, n int) {
	gl.gl.Call("drawArrays", int(mode), i, n)
}
