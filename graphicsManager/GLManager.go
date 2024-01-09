package graphicsManager

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/go-gl/mathgl/mgl32"
)

type GLManager struct {
	Window   *glfw.Window
	Program  uint32
	vaos     []uint32
	vbos     []uint32
	vertices []mgl32.Vec4
	// VertexColors    []mgl32.Vec4
	vec4Storage     Vec4Storage
	float32Storage  Float32Storage
	float32vertices [][]float32
	FS              string
	VS              string
	RenderCall      func()
}

type VerticeStorer interface {
	GetAll() ([][]mgl32.Vec4, error)
	PutSlice([]mgl32.Vec4, string) error
	PutVal(mgl32.Vec4)
}

type Vec4Storage struct {
	ObjectVertices []mgl32.Vec4
	VertexColors   []mgl32.Vec4
}

func (glm *GLManager) NewVec4Storage() Vec4Storage {
	result := Vec4Storage{
		ObjectVertices: []mgl32.Vec4{},
		VertexColors:   []mgl32.Vec4{},
	}

	return result
}

func (glm *GLManager) Vec4Storage() Vec4Storage {
	return glm.vec4Storage
}

func (v4s *Vec4Storage) ClearAll() {
	clear(v4s.ObjectVertices)
	clear(v4s.VertexColors)
}

func (v4s *Vec4Storage) AddOjbVertices(slice []mgl32.Vec4) {
	v4s.ObjectVertices = append(v4s.ObjectVertices, slice...)
}

type Float32Storage struct {
	ObjVecFloats      []float32
	VertexColorFloats []float32
}

func (glm *GLManager) Float32Storage() Float32Storage {
	return glm.float32Storage
}

func (glm *GLManager) NewFloat32Storage() Float32Storage {
	result := Float32Storage{
		ObjVecFloats:      []float32{},
		VertexColorFloats: []float32{},
	}

	return result
}

func (glm *GLManager) GetAll() (slice [][]mgl32.Vec4, err error) {
	for _, val := range slice {
		fmt.Printf("All slices in storage: %v", val)
	}

	return slice, nil
}

func (glm *GLManager) PutSlice(slice []mgl32.Vec4, selection string) (err error) {
	switch selection {
	case "object":
		fmt.Println("Object storage appending")
		fmt.Printf("Value: %v", slice)
		glm.vec4Storage.ObjectVertices = append(glm.vec4Storage.ObjectVertices, slice...)
	case "color":
		fmt.Println("Color storage appending")
		fmt.Printf("Value: %v", slice)
		glm.vec4Storage.VertexColors = append(glm.vec4Storage.VertexColors, slice...)
	default:
		fmt.Println("Must input object or color as selection in string format")
		return fmt.Errorf(err.Error())
	}

	return nil
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

// Any function with Get in front of it is redundant habbit
// while designing. Go best practice is to only create getters for
// struct fields that are private and thus lowercase, allowing for
// a field named vBOs and a get method title VBOs()
func (glm *GLManager) VBOs() []uint32 {
	return glm.vbos
}

func (glm *GLManager) VAOs() []uint32 {
	return glm.vaos
}

func (glm *GLManager) Vertices() []mgl32.Vec4 {
	return glm.vertices
}

func (glm *GLManager) Float32Vertices() [][]float32 {
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

func (glm *GLManager) SetGeoVertices(sliceVec4 []mgl32.Vec4) {
	glm.vec4Storage.ObjectVertices = sliceVec4
	glm.float32Storage.ObjVecFloats = vec4ToFloat32(sliceVec4)
}

func (glm *GLManager) SetColorVertices(sliceVec4 []mgl32.Vec4) {
	glm.vec4Storage.VertexColors = sliceVec4
	glm.float32Storage.VertexColorFloats = vec4ToFloat32(sliceVec4)
}

func (glm *GLManager) GetGeoVertices() []float32 {
	return glm.float32Storage.ObjVecFloats
}

func (glm *GLManager) GetColorVertices() []float32 {
	return glm.float32Storage.VertexColorFloats
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

func (glm *GLManager) BindVAOs() {
	// Multiple Vaos
	for _, vbo := range glm.VBOs() {
		newVao := makeVao(vbo)
		glm.vaos = append(glm.vaos, newVao)
	}
}

// Multiple Vbos
func (glm *GLManager) BindVBOs() {

	newVbo := makeVbo(glm.float32Storage.ObjVecFloats)
	glm.vbos = append(glm.vbos, newVbo)
	newVbo = makeVbo(glm.float32Storage.VertexColorFloats)
	glm.vbos = append(glm.vbos, newVbo)

}

func (glm *GLManager) ConvertVec3ToFloat32() [][]float32 {
	var tempSlice [][]float32

	newSlice := vec4ToFloat32(glm.vertices)
	tempSlice = append(tempSlice, newSlice)

	return tempSlice
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

	return vao
}

func makeVbo(vertices []float32) uint32 {
	fmt.Println("Make vbo called")
	// The first binding of the buffer when called at initialization
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return vbo
}

func vec4ToFloat32(vec4Array []mgl32.Vec4) []float32 {
	// This is my version of "flatten.js" as I am working with mgl32.Vec3 structs in Go to do vector math but need them squashed into an array of float32 to feed the buffer
	float32Array := make([]float32, 0, len(vec4Array)*4)

	for _, vec := range vec4Array {
		float32Array = append(float32Array, vec.X(), vec.Y(), vec.Z(), vec.W())
	}

	return float32Array
}

// RunLoop is where the rendering and buffering take place
func (glm *GLManager) RunLoop(fps int) {
	t := time.Now()
	for !glm.GetWindow().ShouldClose() {

		//Render call
		glm.Render()

		//Check for errors after each call

		glfw.PollEvents()
		glm.GetWindow().SwapBuffers()

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
		t = time.Now()

	}
}
