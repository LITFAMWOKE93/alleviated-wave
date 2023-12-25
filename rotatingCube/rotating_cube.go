package main

import (
	"LITFAMWOKE93/alleviated-wave/graphics_manager"
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	var2DCube = []mgl32.Vec3{
		{0.0, 1.0, 0.0},
		{1.0, 0.0, 0.0},
		{-1.0, 0.0, 0.0},
		{0.0, -1.0, 0.0},
	}

	theta = 0.0
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

	window, err := glfw.CreateWindow(800, 600, "Test Window Instance", nil, nil)
	if err != nil {
		fmt.Println("glfw.CreateWindow() failed:", err)
		return
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		fmt.Println("gl.Init() failed:", err)

	}

	// We need to link a variable in Go to the shader value uTheta

	glm := graphics_manager.GLManager{
		Window: window,
		VS: `
		#version 410

		in vec4 aPosition;
		uniform float uTheta;



		void main() {
    		mat2 rotationMatrix = mat2(cos(uTheta), -sin(uTheta), sin(uTheta), cos(uTheta));
    
    		vec2 rotatedPosition = rotationMatrix * aPosition.xy;

    		gl_Position = vec4(rotatedPosition, 0.0, 1.0);
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

	// Maybe here we send attribute data
	shaderLocName := gl.Str("uTheta" + "\x00")
	// Go strings need to be converted into null-terminated C strings.
	thetaLoc := gl.GetUniformLocation(glm.GetProgram(), shaderLocName)

	// You can send floats, scalars, vectors, matrices to uniform
	glm.BindProgram()

	glm.SetVertices(var2DCube)
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

		// Rotating cube render
		theta += 0.1
		gl.Uniform1f(thetaLoc, float32(theta))
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	}

	gl.Enable(gl.DEPTH_TEST)

	glm.RunLoop(60)

}