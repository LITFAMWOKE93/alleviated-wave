package main

import (
	"LITFAMWOKE93/alleviated-wave/common"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
// a slice of float32, the datatype that is always fed to openGL
// 0, 0, 0 is the center axis of the view

)

const (
	FPS = 60

	// OpenGL needs the termination character \x00 to compile

	// Define the position of the shape, distance of viewport
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

var windowTitle string = "Window Title"

func main() {
	runtime.LockOSThread()

	window := common.InitGlfw("Window Title")

	defer glfw.Terminate()
	program := common.InitOpenGL()

	// Make the object
	t := time.Now()
	for !window.ShouldClose() {

		if err := common.Draw(program, window); err != nil {
			panic(err)
		}
		// create an FPS
		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}

}
