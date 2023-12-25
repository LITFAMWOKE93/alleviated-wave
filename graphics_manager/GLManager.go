package graphics_manager

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/go-gl/mathgl/mgl32"
)

type GLManager struct {
	Window          *glfw.Window
	Program         uint32
	Vao             uint32
	Vbo             uint32
	vertices        []mgl32.Vec3
	float32vertices []float32
	FS              string
	VS              string
	RenderCall      func()
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

func (glm *GLManager) NewProgram() uint32 {

	fmt.Println("OpenGL Version:", gl.GoStr(gl.GetString(gl.VERSION)))

	program, err := newProgram(glm.VertexShaderSource(), glm.FragmentShaderSource())
	if err != nil {
		fmt.Println("Shader program creation failed:", err)
		return 0
	}

	return program

}

//func NewGLManager(width, height int, windowTitle string) GLManager {

//result := GLManager{
//	window: NewWindowContext(width, height, windowTitle),
//TODO: Make setters and getters for uninitilized values
//}

///	return result
//}

func (glm *GLManager) BindProgram() {
	if glm.GetProgram() != 0 {
		gl.UseProgram(glm.GetProgram())
		fmt.Println("BindProgram called")
	} else {
		fmt.Println("Program value is 0")
	}
}

func (glm *GLManager) GetProgram() uint32 {
	return glm.Program
}

func (glm *GLManager) GetWindow() *glfw.Window {
	return glm.Window
}

func (glm *GLManager) VBO() uint32 {
	return glm.Vbo
}

func (glm *GLManager) VAO() uint32 {
	return glm.Vao
}

func (glm *GLManager) Vertices() []mgl32.Vec3 {
	return glm.vertices
}

func (glm *GLManager) Float32Vertices() []float32 {
	return glm.float32vertices
}

func (glm *GLManager) FragmentShaderSource() string {
	return glm.FS
}

func (glm *GLManager) VertexShaderSource() string {
	return glm.VS
}

func (glm *GLManager) SetShaderSource(shaderSource, shaderType string) {

	switch shaderType {
	case "vertex":
		glm.VS = shaderSource
	case "fragment":
		glm.FS = shaderSource
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

func (glm *GLManager) SetProgram() {
	glm.Program = glm.NewProgram()
}

func (glm *GLManager) ClearFloat32Vertices() {
	glm.float32vertices = nil
}

func (glm *GLManager) BindVAO() {
	glm.Vao = makeVao(glm.Vbo)
}

func (glm *GLManager) BindVBO() {
	glm.Vbo = makeVbo(glm.float32vertices)
}

func (glm *GLManager) ConvertVec3ToFloat32() []float32 {
	return vec3ToFloat32(glm.vertices)
}

func (glm *GLManager) Render() {
	if glm.RenderCall != nil {
		glm.RenderCall()
	} else {
		fmt.Println("Render call function is nil")
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {

	// Compile the shaders from the given source, turning it into a uint32 value

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("VertexShaderError")
		return 0, err

	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("FragmentShaderError")
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
	fmt.Println("makeVao called")
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
	fmt.Println("Make vbo called")
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

// RunLoop is where the rendering and buffering take place
func (glm *GLManager) RunLoop(fps int) {
	t := time.Now()
	for !glm.GetWindow().ShouldClose() {
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//Render call
		glm.Render()

		//Check for errors after each call
		for errCode := gl.GetError(); errCode != gl.NO_ERROR; errCode = gl.GetError() {
			fmt.Println("OpenGL error: ", errCode)
		}

		glfw.PollEvents()
		glm.GetWindow().SwapBuffers()
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
		t = time.Now()

	}
}
