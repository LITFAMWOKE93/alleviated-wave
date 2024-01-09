package main

import (
	"fmt"
	"runtime"

	"github.com/LITFAMWOKE93/alleviated-wave/graphicsManager"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Create faces
// Push color information into array
// Interpolate colors across polygon
// Scale the cube in the vertex shader
// Create array of 12 triangles that form 6 faces
// use PRIMITIVE_RESTART_FIXED_INDEX to flag the end of a triangle_Fan
// Value is 255
// Need a fixed point, rotational angle, and vector of rotation
// Move the object to origin, apply rotation, transform back to location
// The transformation specified last is applied first
// Pass as much math as you can to the GPU
// Quaternions prevent gimbal lock

const (
	VERTEXSHADERSOURCE = `
	#version 410

	in vec4 aPosition;
	in vec4 aColor;
	out vec4 vColor;

	uniform vec3 uTheta;

	// quaternion multiplier

	vec4 multq(vec4 a, vec4 b)
	{
		return (vec4(a.x*b.x - dot(a.yzw, b.yzw), a.x*b.yzw+b.x*a.yzq+cross(b.yzw, a.yzw)));
	}

	// inverse quaternion

	vec4 invq(vec4 a)
	{
		return (vec4(a.x, -a.yzw)/dot(a,a));
	}



	void main() {
		vec3 angles = radians( uTheta );
		vec4 r;
		vec4 p;
		vec4 rx, ry, rz;
		vec3 c = cos(angles/2.0);
		vec3 s = sin(angles/2.0);
		rx = vec4(c.x, -s.x, 0.0, 0.0); // x rot quat
		ry = vec4(c.y, 0.0, s.y, 0.0); // y rot quat
		rz = vec4(c.z, 0.0, 0.0, s.z); // z rot quat
		r = multq(rx, multq(ry, rz)); // rot quat
		p = vec4(0.0, aPosition.xyz); // input point quat
		p = multq(r, multq(p, invq(r))); // rotated point quat
		gl_Position = vec4( p.yzw, 1.0); // Convert to homogenous coords
		gl_Position.z = -gl_Position.z; // inverse/reflect
		vColor = aColor;

	}
		` + "\x00"

	FRAGMENTSHADERSOURCE = `
	#version 410
	in vec4 vColor;
	out vec4 fColor;
	void main() {
		fColor = vColor;
	}
		` + "\x00"
)

var (
	numPositions int32 = 36

	Positions []float32
	Colors    []float32

	// Assume the vertices of a cube are available through an array
	// Question: Why do 8 3d vertices construct 6 faces?
	// Answer: The faces are constructed using combinations of the vertices

	v3DCube = []mgl32.Vec4{
		{-0.5, -0.5, 0.5, 1.0},
		{-0.5, 0.5, 0.5, 1.0},
		{0.5, 0.5, 0.5, 1.0},
		{0.5, -0.5, 0.5, 1.0},
		{-0.5, -0.5, -0.5, 1.0},
		{-0.5, 0.5, -0.5, 1.0},
		{0.5, 0.5, -0.5, 1.0},
		{0.5, -0.5, -0.5, 1.0},
	}

	cubeColors = []mgl32.Vec4{
		{0.0, 0.0, 0.0, 1.0}, // black
		{1.0, 0.0, 0.0, 1.0}, // red
		{1.0, 1.0, 0.0, 1.0}, // yellow
		{0.0, 1.0, 0.0, 1.0}, // green
		{0.0, 0.0, 1.0, 1.0}, // blue
		{1.0, 0.0, 1.0, 1.0}, // magenta
		{0.0, 1.0, 1.0, 1.0}, // cyan
		{1.0, 1.0, 1.0, 1.0}, // white
	}
	// Outward facing, vertices traversed in counterclockwise order

	theta          = 0.0
	switchRotation = false
)

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		fmt.Println("glfw.Init() failed:", err)
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Rotating Cube", nil, nil)
	if err != nil {
		fmt.Println("glfw.CreateWindow() failed:", err)
		return
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		fmt.Println("gl.Init() failed:", err)

	}

	// We need to link a variable in Go to the shader value uTheta

	glm := graphicsManager.GLManager{
		Window: window,
		VS:     VERTEXSHADERSOURCE,
		FS:     FRAGMENTSHADERSOURCE,
	}

	glm.NewVec4Storage()
	glm.NewFloat32Storage()

	//glm.Window.SetMouseButtonCallback(mouseEventListener)

	// Set clear color
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.Viewport(0, 0, 800, 600)
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		fmt.Println("OpenGL error after drawing colors:", errCode)
		return
	}

	// Once the shader sources are configured we can create a program

	glm.SetProgram()

	// Maybe here we send attribute data
	shaderLocName := gl.Str("uTheta" + "\x00")
	// Go strings need to be converted into null-terminated C strings.
	thetaLoc := gl.GetUniformLocation(glm.GetProgram(), shaderLocName)

	// You can send floats, scalars, vectors, matrices to uniform
	glm.BindProgram()

	glm.SetGeoVertices(v3DCube)
	glm.SetColorVertices(cubeColors)

	// Create the buffer object that holds the positions, normals, colors and texture coordinates.
	// The vbo can store this data on the GPU
	// Multiple VBO's can be set up
	// TODO: Create a buffer pool and pointers to the last, next, and current buffers for use
	colorCube(glm)

	// Make VBO
	geoVBO := makeVbo(Positions)
	cVBO := makeVbo(Colors)

	gl.BindBuffer(gl.ARRAY_BUFFER, geoVBO)

	fmt.Println("Positions", Positions)
	gl.BufferData(gl.ARRAY_BUFFER, int(4*len(Positions)), gl.Ptr(&Positions[0]), gl.STATIC_DRAW)

	// Find shader variable name
	geoCname := gl.Str("aPosition" + "\x00")
	positionLoc := gl.GetAttribLocation(glm.GetProgram(), geoCname)
	// Make VAO
	geoVAO := makeVao(geoVBO)

	gl.BindVertexArray(geoVAO)
	gl.EnableVertexAttribArray(uint32(positionLoc))

	gl.VertexAttribPointer(uint32(positionLoc), 4, gl.FLOAT, false, 4*4, gl.Ptr(&glm.GetGeoVertices()[0]))

	// Make VBO
	// Bind VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, cVBO)
	// Feed in 32 bytes
	fmt.Println("Colors: ", Colors)
	gl.BufferData(gl.ARRAY_BUFFER, int(4*len(Colors)), gl.Ptr(&Colors[0]), gl.STATIC_DRAW)

	// Find shader variable name
	colorCname := gl.Str("aColor" + "\x00")
	colorLoc := gl.GetAttribLocation(glm.GetProgram(), colorCname)

	if colorLoc == -1 {
		fmt.Println("Failed to find attribute location for aColor")
		return
	}

	// Make Vao
	colorVAO := makeVao(cVBO)

	gl.BindVertexArray(colorVAO)

	gl.EnableVertexArrayAttrib(colorVAO, uint32(colorLoc))

	gl.VertexAttribPointer(uint32(colorLoc), 4, gl.FLOAT, false, 4*4, gl.Ptr(&glm.GetColorVertices()[0]))

	//glm.BindVBOs()
	//glm.BindVAOs()

	glm.RenderCall = func() {

		if errCode := gl.GetError(); errCode != gl.NO_ERROR {
			fmt.Println("OpenGL error before rendering:", errCode)
			return
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Rotating cube render

		theta += 2.0

		// Convert theta to float32 and create a slice
		thetaFloat32 := float32(theta)
		thetaSlice := []float32{thetaFloat32, 0.0, 0.0}

		// Update the uniform
		gl.Uniform3fv(thetaLoc, 1, &thetaSlice[0])
		gl.BindVertexArray(geoVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, numPositions)

		if errCode := gl.GetError(); errCode != gl.NO_ERROR {
			fmt.Println("OpenGL error after drawing geometry:", errCode)
			return
		}

		gl.BindVertexArray(colorVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, numPositions)

		if errCode := gl.GetError(); errCode != gl.NO_ERROR {
			fmt.Println("OpenGL error after drawing colors:", errCode)
			return
		}

	}

	glm.RunLoop(60)

}

// VERY BAD: Using a language built with the purpose of composition to basically couple it with inheritance

func colorCube(glm graphicsManager.GLManager) {
	Quad(1, 0, 3, 2, glm)
	Quad(2, 3, 7, 6, glm)
	Quad(3, 0, 4, 7, glm)
	Quad(6, 5, 1, 2, glm)
	Quad(4, 5, 6, 7, glm)
	Quad(5, 4, 0, 1, glm)

}

func makeVbo(vertices []float32) uint32 {

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return vbo
}

func makeVao(vbo uint32) uint32 {

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	return vao
}

func Quad(a, b, c, d int, glm graphicsManager.GLManager) {

	vertices := vec4ToFloat32(glm.Vec4Storage().ObjectVertices)

	vertexColors := vec4ToFloat32(glm.Vec4Storage().VertexColors)

	var indices = []int{a, b, c, a, c, d}

	for _, val := range indices {

		Positions = append(Positions, vertices[val])

		Colors = append(Colors, vertexColors[val])

	}
}

func vec4ToFloat32(vec4Array []mgl32.Vec4) []float32 {
	// This is my version of "flatten.js" as I am working with mgl32.Vec3 structs in Go to do vector math but need them squashed into an array of float32 to feed the buffer
	float32Array := make([]float32, 0, len(vec4Array)*4)

	for _, vec := range vec4Array {
		float32Array = append(float32Array, vec.X(), vec.Y(), vec.Z(), vec.W())
	}

	return float32Array
}
