package common

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	WIDTH  = 1280
	HEIGHT = 800

	vertexShaderSource = `
		#version 400

		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"
	// Define the color of the shape
	fragmentShaderSource = `
		#version 400

		out vec4 frag_colour;
		void main() {
  			frag_colour = vec4(1, 1, 1, 0.3);
		}
	` + "\x00"
)

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(string windowName) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	wN := windowName

	window, err := glfw.CreateWindow(WIDTH, HEIGHT, wN, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and return an init program.
func initOpenGL() uint32 {

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	// Attaching the shaders to the shader program
	gl.AttachShader(prog, vShader)
	gl.AttachShader(prog, fShader)

	// Linking program to the GL context
	gl.LinkProgram(prog)
	return prog
}

// makeVertexArrayObject inits and returns a vertex array from the points provided.
func makeVertexArrayObject(points []float32) uint32 {
	var vertexBufferObject uint32

	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)

	// a 32-bit float has 4 bytes, so we are saying the size of teh buffer in bytes is 4 times the number of points provided.
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vertexArrayObject

}

// func draw redraws the window everything in the frame
func draw(prog uint32, window *glfw.Window) error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	glfw.PollEvents()
	window.SwapBuffers()
	return nil

}

// compileShader recieves the shader source as a string and it's type and returns a pointer
// to the resulting compiled shader, log message on failure to compile
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	if source != "" {
		csources, free := gl.Strs(source)

		gl.ShaderSource(shader, 1, csources, nil)
		free()
		gl.CompileShader(shader)

	}

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))

		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}
