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
		retun (vec4(a.x*b.x - dot(a.yzw, b.yzw), a.x*b.yzw+b.x*a.yzq+cross(b.yzw, a.yzw)));
	}

	// inverse quaternion

	vec4 invq(vec4 a)
	{
		return (vec4(a.x, -a.yzw)/dot(a,a));
	}



	void main() {
		vec3 angle = radians( uTheta );
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
	glm.Window.SetMouseButtonCallback(mouseEventListener)

	// Set clear color
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// Set shader sources

	fmt.Println(glm.FragmentShaderSource())
	fmt.Println(glm.VertexShaderSource())

	// Once the shader sources are configured we can create a program

	glm.SetProgram()

	// Maybe here we send attribute data
	shaderLocName := gl.Str("uTheta" + "\x00")
	// Go strings need to be converted into null-terminated C strings.
	thetaLoc := gl.GetUniformLocation(glm.GetProgram(), shaderLocName)

	// You can send floats, scalars, vectors, matrices to uniform
	glm.BindProgram()

	fmt.Println("Instance vec3 slice:", glm.Vertices())

	glm.SetVertices(v3DCube)
	glm.SetVertices(cubeColors)
	glm.SetFloat32Vertices()
	fmt.Println("Instance float32 vertices", glm.Float32Vertices())
	// Create the buffer object that holds the positions, normals, colors and texture coordinates.
	// The vbo can store this data on the GPU
	// Multiple VBO's can be set up
	// TODO: Create a buffer pool and pointers to the last, next, and current buffers for use

	glm.BindVBOs()
	glm.BindVAOs()

	glm.RenderCall = func() {

		// Rotating cube render
		if !switchRotation {
			theta += 0.1
		} else {
			theta -= 0.1
		}
		gl.Uniform1f(thetaLoc, float32(theta))

		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, numPositions)
		// Render button

	}

	gl.Enable(gl.DEPTH_TEST)

	glm.RunLoop(60)

}

func mouseEventListener(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		// Check mousePos in button area
		x, y := w.GetCursorPos()
		if x >= 640 && x <= 800 && y >= 480 && y <= 600 {
			switchRotation = !switchRotation
		}
	}
}

// instead of pushing to global arrays quad hands the vertices to the manager
func quad(a, b, c, d int) {

}
