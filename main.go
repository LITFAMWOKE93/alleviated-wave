package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	cubeVertices = []mgl32.Vec3{
		{-0.5, -0.5, -0.5},
		{0.5, -0.5, -0.5},
		{0.5, 0.5, -0.5},
		{-0.5, 0.5, -0.5},
		{-0.5, -0.5, 0.5},
		{0.5, -0.5, 0.5},
		{0.5, 0.5, 0.5},
		{-0.5, 0.5, 0.5},
	}

	cubeIndices = []uint32{
		0, 1, 2, 2, 3, 0, // Front face
		4, 5, 6, 6, 7, 4, // Back face
		0, 3, 7, 7, 4, 0, // Left face
		1, 2, 6, 6, 5, 1, // Right face
		3, 2, 6, 6, 7, 3, // Top face
		0, 1, 5, 5, 4, 0, // Bottom face
	}
)

func main() {
	runtime.LockOSThread()

	//TODO: Create a GL struct that handles the gl library through composition

	if err := glfw.Init(); err != nil {
		fmt.Println("glfw.Init() failed:", err)
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Test Window Instance", nil, nil)
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
		VS: `
		#version 410
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
			` + "\x00",
		FS: `
		#version 410
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1.0, 0.0, 0.0, 1.0);
		}
			` + "\x00",
	}

	// Set shader sources

	fmt.Println(glm.FragmentShaderSource())
	fmt.Println(glm.VertexShaderSource())

	// Once the shader sources are configured we can create a program

	glm.SetProgram()
	glm.BindProgram()

	glm.SetVertices(cubeVertices)
	fmt.Println("Instance vec3 slice:", glm.Vertices())

	glm.SetFloat32Vertices()
	fmt.Println("Instance float32 vertices", glm.Float32Vertices())
	// Create the buffer object that holds the positions, normals, colors and texture coordinates.
	// The vbo can store this data on the GPU
	// Multiple VBO's can be set up
	// TODO: Create a buffer pool and pointers to the last, next, and current buffers for use

	glm.BindVBO()
	fmt.Println("Instance VBO: ", glm.VBO())

	glm.BindVAO()
	fmt.Println("Instance VAO: ", glm.VAO())

	glm.RenderCall = func() {

		// Drawing for cube
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

		fmt.Println("Render call")
		//fmt.Println("VAO", glm.VAO())
		//fmt.Println("VBO", glm.VBO())

	}

	gl.Enable(gl.DEPTH_TEST)

	glm.RunLoop(30)

}
