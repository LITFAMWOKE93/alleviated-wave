package main

import (
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	// a slice of float32, the datatype that is always fed to openGL
	// 0, 0, 0 is the center axis of the view
	vertices []common.Vec32
)

const (
	// Frames per second for sleep timer
	FPS = 60
	// How many points to generate
	numPositions = 5000
)

func init() {
	// The algorithm for generating the gasket
	// Hacky, change later. Load the vertices into an array
	var vec1 common.Vec32 = common.NewVec32(-1.0, -1.0, 0.0)
	var vec2 common.Vec32 = common.NewVec32(0.0, 1.0, 0.0)
	var vec3 common.Vec32 = common.NewVec32(1.0, -1.0, 0.0)

	vertices = append(vertices, vec1)
	vertices = append(vertices, vec2)
	vertices = append(vertices, vec3)

	var u = common.Vec32.Add()

}

func main() {
	runtime.LockOSThread()

	window := common.InitGlfw("Sierpenski Gasket")

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
