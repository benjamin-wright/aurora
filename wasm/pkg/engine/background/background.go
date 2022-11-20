package background

import (
	"fmt"

	"aurora.ponglehub.co.uk/pkg/engine/fetch"
	"aurora.ponglehub.co.uk/pkg/webgl"
)

type Background struct {
	Y float32
	X float32

	texture []uint8
	normals []float32

	meshBuffer     webgl.Buffer
	elementBuffer  webgl.Buffer
	texCoordBuffer webgl.Buffer
	textureObj     webgl.Texture
	texBuffer      webgl.Buffer
	normalBuffer   webgl.Buffer
}

func New(image string) (*Background, error) {
	texture, err := fetch.Png(fmt.Sprintf("content/main/backgrounds/%s.png", image))
	if err != nil {
		return nil, fmt.Errorf("failed to get texture file: %+v", err)
	}

	normals, err := fetch.Float32(fmt.Sprintf("content/main/backgrounds/%s.normals", image))
	if err != nil {
		return nil, fmt.Errorf("failed to get normals file: %+v", err)
	}

	return &Background{
		texture: texture,
		normals: normals,
	}, nil
}

func (b *Background) Init(gl *webgl.WebGL, program webgl.Program) {
	b.meshBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, b.meshBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, webgl.Float32ArrayBuffer(backgroundMesh), gl.STATIC_DRAW)

	// b.elementBuffer = gl.CreateBuffer()
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.elementBuffer)
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, webgl.Uint16ArrayBuffer(backgroundIndices), gl.STATIC_DRAW)

	b.textureObj = gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, &b.textureObj)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 256, 256, 0, gl.RGBA, gl.UNSIGNED_BYTE, b.texture)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	b.texCoordBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, b.texCoordBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, webgl.Float32ArrayBuffer(backgroundTexCoords), gl.STATIC_DRAW)
}

func (b *Background) Render(gl *webgl.WebGL, program webgl.Program) {
	gl.UseProgram(program)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.meshBuffer)

	positionAttrib := gl.GetAttribLocation(program, "position")
	gl.VertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(positionAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.texCoordBuffer)

	textureAttrib := gl.GetAttribLocation(program, "texture")
	gl.VertexAttribPointer(textureAttrib, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(textureAttrib)

	// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_SHORT, 0)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

const backgroundVertexShader = `#version 300 es
in vec2 position;
in vec2 texture;

out highp vec2 vTexture;

void main(void) {
    gl_Position = vec4(position, 0.0, 1.0) * vec4(1.0, -1.0, 1.0, 1.0);
    vTexture = texture;
}
`

const backgroundFragmentShader = `#version 300 es
in highp vec2 vTexture;
uniform sampler2D uSampler;

out mediump vec4 FragColor;

void main(void) {
	FragColor = texture(uSampler, vTexture);
}
`

var backgroundMesh = []float32{
	-0.8, -0.8,
	-0.8, 0.8,
	0.8, -0.8,
	0.8, -0.8,
	-0.8, 0.8,
	0.8, 0.8,
}

var backgroundIndices = []uint16{
	0, 1, 2,
	2, 1, 3,
}

var backgroundTexCoords = []float32{
	0, 0,
	0, 1,
	1, 0,
	1, 0,
	0, 1,
	1, 1,
}

func CompileProgram(gl *webgl.WebGL) (webgl.Program, error) {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, backgroundVertexShader)
	gl.CompileShader(vertexShader)
	if !gl.GetShaderParameter(vertexShader, gl.COMPILE_STATUS).(bool) {
		compilationLog := gl.GetShaderInfoLog(vertexShader)
		return webgl.Program{}, fmt.Errorf("compile failed (VERTEX_SHADER) %v", compilationLog)
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, backgroundFragmentShader)
	gl.CompileShader(fragmentShader)
	if !gl.GetShaderParameter(fragmentShader, gl.COMPILE_STATUS).(bool) {
		compilationLog := gl.GetShaderInfoLog(fragmentShader)
		return webgl.Program{}, fmt.Errorf("compile failed (FRAGMENT_SHADER) %v", compilationLog)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	if !gl.GetProgramParameter(program, gl.LINK_STATUS).(bool) {
		return webgl.Program{}, fmt.Errorf("link failed: %v", gl.GetProgramInfoLog(program))
	}

	return program, nil
}
