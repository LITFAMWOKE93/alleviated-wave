package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
	FPS    = 60
)

var (
	vertices = []mgl32.Vec3{
		{0.0, 1.0, 0.0},
		{-1.0, -1.0, 0.0},
		{1.0, -1.0, 0.0},
	}

	VAO uint32
	VBO uint32

	float32vertices = vec3ToFloat32(vertices)

	vertexShaderSource = `
		#version 410
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1.0, 0.0, 0.0, 1.0);
		}
	` + "\x00"
)

//Conceptually I thought I would be generating the points one frame at a time but now realize that the
//vertice calculations were done up front and handed to the buffer so that it could render a gasket in a single frame

func main() {
	runtime.LockOSThread()

	fmt.Println("Float32Vert value at init: ", float32vertices)

	if err := glfw.Init(); err != nil {
		fmt.Println("glfw.Init() failed:", err)
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Sierpinski's Gasket", nil, nil)
	if err != nil {
		fmt.Println("glfw.CreateWindow() failed:", err)
		return
	}
	window.MakeContextCurrent()

	window.SetKeyCallback(keyCallback)

	if err := gl.Init(); err != nil {
		fmt.Println("gl.Init() failed:", err)
		return
	}

	VBO = makeVbo(float32vertices)
	VAO = makeVao(VBO)

	gl.Enable(gl.DEPTH_TEST)

	fmt.Println("OpenGL version:", gl.GoStr(gl.GetString(gl.VERSION)))

	program, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		fmt.Println("Shader program creation failed:", err)
		return
	}

	for errCode := gl.GetError(); errCode != gl.NO_ERROR; errCode = gl.GetError() {
		fmt.Println("OpenGL error: ", errCode)
	}
	gl.UseProgram(program)
	// TODO: turn this into a draw function

	t := time.Now()
	for !window.ShouldClose() {
		// Call ClearColor before color
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.BindVertexArray(VAO)

		// Render call

		// The main loop handles the actual draw calls and parsing of the buffer data into the
		// the correct format for reading into buffer

		for errCode := gl.GetError(); errCode != gl.NO_ERROR; errCode = gl.GetError() {
			fmt.Println("OpenGL error: ", errCode)
		}
		float32vertices = nil
		glfw.PollEvents()
		window.SwapBuffers()
		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
		t = time.Now()

	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {

	fmt.Println("Vertex Shader Source:\n", vertexShaderSource)
	fmt.Println("Fragment Shader Source:\n", fragmentShaderSource)
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetProgramInfoLog(program, logLength, nil, &log[0])

		return 0, fmt.Errorf("program link error: %s", string(log))
	}

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])

		return 0, fmt.Errorf("shader compile error: %s", string(log))
	}

	return shader, nil
}

func vec3ToFloat32(vec3Array []mgl32.Vec3) []float32 {
	float32Array := make([]float32, 0, len(vec3Array)*3)

	for _, vec := range vec3Array {
		float32Array = append(float32Array, vec.X(), vec.Y(), vec.Z())
	}

	return float32Array
}

func makeVao(vbo uint32) uint32 {

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	// I was generating an empty buffer here, 5 hours to find.

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(0)

	return vao
}

func renderGasket(v0, v1, v2 mgl32.Vec3, depth int) {
	if depth == 0 {
		float32vertices = pushTriangle(float32vertices, v0, v1, v2)
		return
	}

	// Calculate midpoints of edges
	mid01 := v0.Add(v1).Mul(0.5)
	mid12 := v1.Add(v2).Mul(0.5)
	mid20 := v2.Add(v0).Mul(0.5)

	//fmt.Printf("Depth: %d, Vertices: (%v, %v, %v)\n", depth, v0, v1, v2)

	// Recursive calls for three sub-triangles
	renderGasket(v0, mid01, mid20, depth-1)
	renderGasket(mid01, v1, mid12, depth-1)
	renderGasket(mid20, mid12, v2, depth-1)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(float32vertices), gl.Ptr(float32vertices), gl.STATIC_DRAW)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(float32vertices)/3))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

}

func pushTriangle(vertices []float32, v0, v1, v2 mgl32.Vec3) []float32 {
	return append(vertices, v0.X(), v0.Y(), v0.Z(), v1.X(), v1.Y(), v1.Z(), v2.X(), v2.Y(), v2.Z())
}
func makeVbo(vertices []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return vbo
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// Check if the key is pressed and not released

	if action == glfw.Press || action == glfw.Repeat {
		switch key {
		case glfw.Key1:
			// Render when '1' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 1)
			fmt.Println("Key 1 pressed")
		case glfw.Key2:
			// Render when '1' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 2)
			fmt.Println("Key 2 pressed")
		case glfw.Key3:
			// Render when '1' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 3)
			fmt.Println("Key 3 pressed")
		case glfw.Key4:
			// Render when '1' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 4)
			fmt.Println("Key 4 pressed")
		case glfw.Key5:
			// Render when '1' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 5)
			fmt.Println("Key 5 pressed")
		case glfw.Key6:
			// Render when '1' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 6)
			fmt.Println("Key 6 pressed")
		}

	}
}
