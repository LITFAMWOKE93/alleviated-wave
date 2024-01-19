package main

import (
	"fmt"
	"math"
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
		return (vec4(a.x*b.x - dot(a.yzw, b.yzw), a.x*b.yzw+b.x*a.yzw+cross(b.yzw, a.yzw)));
	}

	// inverse quaternion

	vec4 invq(vec4 a)
	{
		return (vec4(a.x, -a.yzw)/dot(a,a));
	}

	uniform mat4 uModelViewMatrix;
	uniform mat4 uProjectionMatrix;



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
		vec4 utransformedPosition = vec4(p.yzw, 1.0); // Convert to homogeneous coordinates
		
		vColor = aColor;
		gl_Position = uProjectionMatrix * uModelViewMatrix * aPosition;
		gl_Position.z = -gl_Position.z; // inverse/reflect

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

const (
	X_AXIS = 0
	Y_AXIS = 1
	Z_AXIS = 2
)

var (
	numPositions int32 = 36

	at        = mgl32.Vec4{0.0, 0.0, 0.0, 0.0}
	up        = mgl32.Vec4{0.0, 1.0, 0.0, 0.0}
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

	theta  = []float32{0.0, 0.0, 0.0}
	near   = float32(-1)
	far    = float32(1)
	radius = float32(1)

	phi = 0.0

	left   = float32(-1.0)
	right  = float32(1.0)
	bottom = float32(-1.0)
	top    = float32(1.0)

	//Event handler variables

	isMousePressed         bool
	lastMouseX, lastMouseY float64
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

	glm := graphicsManager.GLManager{
		Window: window,
		VS:     VERTEXSHADERSOURCE,
		FS:     FRAGMENTSHADERSOURCE,
	}

	glm.NewVec4Storage()
	glm.NewFloat32Storage()

	glm.Window.SetMouseButtonCallback(mouseEventListener)

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
	// The SINGLE vbo can store this data on the GPU for each unique object
	// Multiple VBO's can be set up
	// TODO: Create a buffer pool and pointers to the last, next, and current buffers for use
	colorCube(glm)

	// Find shader variable name
	geoCname := gl.Str("aPosition" + "\x00")
	positionLoc := gl.GetAttribLocation(glm.GetProgram(), geoCname)

	colorCname := gl.Str("aColor" + "\x00")
	colorLoc := gl.GetAttribLocation(glm.GetProgram(), colorCname)

	modelViewCname := gl.Str("uModelViewMatrix" + "\x00")
	modelViewMatLoc := gl.GetUniformLocation(glm.GetProgram(), modelViewCname)

	projMatCname := gl.Str("uProjectionMatrix" + "\x00")
	projMatLoc := gl.GetUniformLocation(glm.GetProgram(), projMatCname)
	// Bind VAO.
	// Bind VBO.
	// Set vertex attribute pointers.
	// Unbind VBO
	// Unbind VAO.
	geoVBO := makeVbo()
	cVBO := makeVbo()
	VAO := makeVao()

	// Create and bind a single VAO
	gl.BindVertexArray(VAO)

	// Setup for position data
	gl.BindBuffer(gl.ARRAY_BUFFER, geoVBO)
	gl.BufferData(gl.ARRAY_BUFFER, int(len(Positions)*4*4), gl.Ptr(&Positions[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(uint32(positionLoc))
	gl.VertexAttribPointer(uint32(positionLoc), 4, gl.FLOAT, false, 0, nil)

	// Setup for color data
	gl.BindBuffer(gl.ARRAY_BUFFER, cVBO)
	gl.BufferData(gl.ARRAY_BUFFER, int(len(Colors)*4*4), gl.Ptr(&Colors[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(uint32(colorLoc))
	gl.VertexAttribPointer(uint32(colorLoc), 4, gl.FLOAT, false, 0, nil)

	// Unbind VBO and VAO
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	glm.RenderCall = func() {
		if errCode := gl.GetError(); errCode != gl.NO_ERROR {
			fmt.Println("OpenGL error before rendering:", errCode)
			return
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// Create polar coordinates for the eye, when looking at the origin of object coordinates
		eye := mgl32.Vec4{radius * float32(math.Sin(float64(theta[0]))) * float32(math.Cos(float64(phi))),
			radius * float32(math.Sin(float64(theta[0]))) * float32(math.Sin(float64(phi))),
			radius * float32(math.Cos(float64(theta[0]))), 1.0}
		// Create the model view matrix using the u v n properties, looking at the origin
		modelViewMatrix := mgl32.LookAt(eye.X(), eye.Y(), eye.Z(), at.X(), at.Y(), at.Z(), up.X(), up.Y(), up.Z())

		gl.UniformMatrix4fv(modelViewMatLoc, 1, false, &modelViewMatrix[0])
		// An orthographic projection
		projectionMatrix := ortho(left, right, bottom, top, near, far)
		// Give the information to the Shader
		gl.UniformMatrix4fv(projMatLoc, 1, false, &projectionMatrix[0])
		// Rotating cube render
		updateRotation(glm.Window)

		// Update the uniform
		gl.Uniform3fv(thetaLoc, 1, &theta[0])

		// Bind the single VAO
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, numPositions)

		if errCode := gl.GetError(); errCode != gl.NO_ERROR {
			fmt.Println("OpenGL error after drawing:", errCode)
			return
		}
	}

	glm.RunLoop(60)

}

func colorCube(glm graphicsManager.GLManager) {
	Quad(1, 0, 3, 2, glm)
	Quad(2, 3, 7, 6, glm)
	Quad(3, 0, 4, 7, glm)
	Quad(6, 5, 1, 2, glm)
	Quad(4, 5, 6, 7, glm)
	Quad(5, 4, 0, 1, glm)

}

func makeVbo() uint32 {

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return vbo
}

func makeVao() uint32 {

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	return vao
}

func Quad(a, b, c, d int, glm graphicsManager.GLManager) {

	vertices := glm.Vec4Storage().ObjectVertices

	vertexColors := glm.Vec4Storage().VertexColors

	var indices = []int{a, b, c, a, c, d}

	for i := 0; i < len(indices); i++ {
		fmt.Println(i)
		Positions = append(Positions, vecToFloat32(vertices[indices[i]])...)

		Colors = append(Colors, vecToFloat32(vertexColors[a])...)
	}

}

func vecToFloat32(vec mgl32.Vec4) []float32 {
	float32Array := make([]float32, 0, 4)

	float32Array = append(float32Array, vec.X(), vec.Y(), vec.Z(), vec.W())
	return float32Array
}

func mouseEventListener(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			isMousePressed = true
			lastMouseX, lastMouseY = window.GetCursorPos()
		} else if action == glfw.Release {
			isMousePressed = false
		}
	}
}

func updateRotation(window *glfw.Window) {
	if isMousePressed {
		x, y := window.GetCursorPos()
		deltaX := x - lastMouseX
		deltaY := y - lastMouseY

		// Update theta here based on deltaX and deltaY
		// The sensitivity factor controls how much the rotation changes with mouse movement
		var sensitivity float32 = 0.1
		theta[0] += float32(deltaY) * sensitivity
		theta[1] -= float32(deltaX) * sensitivity

		// Update last mouse position
		lastMouseX = x
		lastMouseY = y
	}
}

// Performs the scalar transformation
func ortho(left, right, bottom, top, near, far float32) mgl32.Mat4 {

	switch {
	case left == right:
		fmt.Println(" Ortho(): Left and right are equal")
	case bottom == top:
		fmt.Println(" Ortho(): bottom and top are equal")
	case near == far:
		fmt.Println(" Ortho(): Near and far are equal")
	}

	width := right - left
	height := top - bottom
	depth := far - near

	result := mgl32.Mat4{}
	result.Set(0, 0, (2.0 / width))
	result.Set(1, 1, (2.0 / height))
	result.Set(2, 2, (2.0 / depth))

	result.Set(0, 3, -(left+right)/width)
	result.Set(1, 3, -(top+bottom)/height)
	result.Set(2, 3, -(near+far)/depth)

	result.Set(3, 3, 1.0)

	return result

}
