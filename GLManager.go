package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/go-gl/mathgl/mgl32"
)

type GLManager struct {
	window          *glfw.Window
	program         uint32
	vao             uint32
	vbo             uint32
	vertices        []mgl32.Vec3
	float32vertices []float32
	fS              string
	vS              string
}

func NewWindowContext(width, height int, windowTitle string) *glfw.Window {

	if err := glfw.Init(); err != nil {
		fmt.Println("glfw.Init() failed:", err)
		return nil
	}
	defer glfw.Terminate()

	// Information needed for the window

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Creating of the window context

	window, err := glfw.CreateWindow(width, height, windowTitle, nil, nil)
	if err != nil {
		fmt.Println("glfw.CreateWindow() failed:", err)
		return nil
	}

	// Binding of the window context
	window.MakeContextCurrent()

	return window

}

func NewProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	if err := gl.Init(); err != nil {
		fmt.Println("gl.Init() failed:", err)
		return 0
	}

	program, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		fmt.Println("Shader program creation failed:", err)
		return 0
	}

	return program

}

func NewGLManager(width, height int, windowTitle string) GLManager {

	result := GLManager{
		window: NewWindowContext(width, height, windowTitle),
		//TODO: Make setters and getters for uninitilized values
	}

	return result
}

func (glm *GLManager) BindProgram() {
	if glm.Program() != 0 {
		gl.UseProgram(glm.Program())
	}
}

func (glm *GLManager) Program() uint32 {
	return glm.program
}

func (glm *GLManager) Window() *glfw.Window {
	return glm.window
}

func (glm *GLManager) VBO() uint32 {
	return glm.vbo
}

func (glm *GLManager) VAO() uint32 {
	return glm.vao
}

func (glm *GLManager) Vertices() []mgl32.Vec3 {
	return glm.vertices
}

func (glm *GLManager) FragmentShaderSource() string {
	return glm.fS
}

func (glm *GLManager) VertexShaderSource() string {
	return glm.vS
}

func (glm *GLManager) SetShaderSrouce(shaderSource, shaderType string) {

	switch shaderType {
	case "vertex":
		glm.vS = shaderSource
	case "fragment":
		glm.fS = shaderSource
	default:
		fmt.Println("Unsupported shader type, please declare \"vertex\" or \"fragment\"")
	}

}

func (glm *GLManager) SetVertices(sliceVec3 []mgl32.Vec3) {
	glm.vertices = sliceVec3
}

func (glm *GLManager) ClearVertices() {
	glm.vertices = nil
}

func (glm *GLManager) SetFloat32Vertices() {
	if glm.Vertices() != nil {
		glm.float32vertices = glm.ConvertVec3ToFloat32()
	}

}

func (glm *GLManager) ClearFloat32Vertices() {
	glm.float32vertices = nil
}

func (glm *GLManager) BindVAO() {
	glm.vao = makeVao(glm.vbo)
}

func (glm *GLManager) BindVBO() {
	float32array := glm.ConvertVec3ToFloat32()
	glm.vbo = makeVbo(float32array)
}

func (glm *GLManager) ConvertVec3ToFloat32() []float32 {
	return vec3ToFloat32(glm.vertices)
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

func vec3ToFloat32(vec3Array []mgl32.Vec3) []float32 {
	// This is my version of "flatten.js" as I am working with mgl32.Vec3 structs in Go to do vector math but need them squashed into an array of float32 to feed the buffer
	float32Array := make([]float32, 0, len(vec3Array)*3)

	for _, vec := range vec3Array {
		float32Array = append(float32Array, vec.X(), vec.Y(), vec.Z())
	}

	return float32Array
}
