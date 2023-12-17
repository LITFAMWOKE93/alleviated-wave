package main

import (
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
// a slice of float32, the datatype that is always fed to openGL
// 0, 0, 0 is the center axis of the view

)

const (
	WIDTH  = 800
	HEIGHT = 600

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

func main() {
	runtime.LockOSThread()

	window := initGlfw("Window Title")

	defer glfw.Terminate()
	program := initOpenGL()

	// Make the object

	for !window.ShouldClose() {

		t := time.Now()
		for x := range cells {
			for _, c := range cells[x] {
				c.checkState(cells)
			}
		}
		if err := draw(program, window, cells); err != nil {
			panic(err)
		}

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}

}
