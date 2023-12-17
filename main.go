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
)

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
