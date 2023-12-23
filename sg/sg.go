package sg

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

	//TODO: Turn into a struct that holds information related to the window, openGL program, and vertex and buffer information.

	// Window initialization using the gl-go/glfw package which acts as the glue for the OS

	if err := glfw.Init(); err != nil {
		fmt.Println("glfw.Init() failed:", err)
		return
	}
	defer glfw.Terminate()

	// Information needed for the window

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Creating of the window context

	window, err := glfw.CreateWindow(width, height, "Sierpinski's Gasket", nil, nil)
	if err != nil {
		fmt.Println("glfw.CreateWindow() failed:", err)
		return
	}

	// Binding of the window context
	window.MakeContextCurrent()

	window.SetKeyCallback(keyCallback)

	if err := gl.Init(); err != nil {
		fmt.Println("gl.Init() failed:", err)
		return
	}

	// Vertex buffer object and vertex binder attribute are initialized and bound before rendering.

	VBO = makeVbo(float32vertices)
	VAO = makeVao(VBO)

	// Depth testing is important for rendering 2D objects in 3D space, checking for vertexes clipping one another

	gl.Enable(gl.DEPTH_TEST)

	fmt.Println("OpenGL version:", gl.GoStr(gl.GetString(gl.VERSION)))

	// Load the shader sources into the program

	program, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		fmt.Println("Shader program creation failed:", err)
		return
	}

	for errCode := gl.GetError(); errCode != gl.NO_ERROR; errCode = gl.GetError() {
		fmt.Println("OpenGL error: ", errCode)
	}

	// Bind the program
	gl.UseProgram(program)
	// TODO: turn this into a draw function

	// Time is used to create a frame per second display to avoid the rendering from happening too quickly and closing the window
	t := time.Now()
	for !window.ShouldClose() {
		// Call ClearColor before Clear, sets the color that will be used to clear the screen
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render call
		renderGasket(vertices[0], vertices[1], vertices[2], 6)

		// The main loop handles the actual draw calls and parsing of the buffer data into the
		// the correct format for reading into buffer

		for errCode := gl.GetError(); errCode != gl.NO_ERROR; errCode = gl.GetError() {
			fmt.Println("OpenGL error: ", errCode)
		}

		// The vertices are reset every loop to avoid appending infinitely
		float32vertices = nil
		// PollEvents is an event listener looking for keybind interactions/mouse clicks
		glfw.PollEvents()
		// SwapBuffers is used to double buffer the frame, to make sure the frame is only visible when fully rendered
		// will see a blank screen without a double buffer
		window.SwapBuffers()
		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
		t = time.Now()

	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {

	fmt.Println("Vertex Shader Source:\n", vertexShaderSource)
	fmt.Println("Fragment Shader Source:\n", fragmentShaderSource)

	// Compile the shaders from the given source, turning it into a uint32 value
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	// Create program, attach shaders, and link them
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	// Check program attributes for errors
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
	// Create a shader uint32 from the c sources provided
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
	// This is my version of "flatten.js" as I am working with mgl32.Vec3 structs in Go to do vector math but need them squashed into an array of float32 to feed the buffer
	float32Array := make([]float32, 0, len(vec3Array)*3)

	for _, vec := range vec3Array {
		float32Array = append(float32Array, vec.X(), vec.Y(), vec.Z())
	}

	return float32Array
}

func makeVao(vbo uint32) uint32 {
	// Vertex array is generated and bound, attributes are set up so we know how many to read at a time and what data type will be read

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	// I was generating an empty buffer here, 5 hours to find.

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(0)

	return vao
}

// renderGasket is the recursive vector math to produce the fractal image, returning the values to the buffer and feeding once the recursion is complete
func renderGasket(v0, v1, v2 mgl32.Vec3, depth int) {
	// The recursive call for the fractal rendering

	//
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

	// Update buffer bindings for the set of float32vertices for rendering
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(float32vertices), gl.Ptr(float32vertices), gl.STATIC_DRAW)

	// The draw call using triangle primitives
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(float32vertices)/3))
	//gl.DrawArrays(gl.LINES, 0, int32(len(float32vertices)/3))
	// Using the POINTS primitive will only render the dot location of each vertice instead of connecting them like the triangle primitive
	// gl.DrawArrays(gl.POINTS, 0, int32(len(float32vertices)/3))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

}

func pushTriangle(vertices []float32, v0, v1, v2 mgl32.Vec3) []float32 {
	// Take the indiviual float32 values and append them
	return append(vertices, v0.X(), v0.Y(), v0.Z(), v1.X(), v1.Y(), v1.Z(), v2.X(), v2.Y(), v2.Z())
}
func makeVbo(vertices []float32) uint32 {
	// The first binding of the buffer when called at initialization
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// 32 bits 4 bytes
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return vbo
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// Check if the key is pressed and not released
	//Keyboard interaction to render different depths
	if action == glfw.Press || action == glfw.Repeat {
		switch key {
		case glfw.Key1:
			// Render when '1' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 1)
			fmt.Println("Key 1 pressed")
		case glfw.Key2:
			// Render when '2' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 2)
			fmt.Println("Key 2 pressed")
		case glfw.Key3:
			// Render when '3' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 3)
			fmt.Println("Key 3 pressed")
		case glfw.Key4:
			// Render when '4' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 4)
			fmt.Println("Key 4 pressed")
		case glfw.Key5:
			// Render when '5' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 5)
			fmt.Println("Key 5 pressed")
		case glfw.Key6:
			// Render when '6' key is pressed
			renderGasket(vertices[0], vertices[1], vertices[2], 6)
			fmt.Println("Key 6 pressed")
		}

	}
}
