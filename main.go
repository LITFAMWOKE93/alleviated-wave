package main

import (
	"runtime"
)

var (
// a slice of float32, the datatype that is always fed to openGL
// 0, 0, 0 is the center axis of the view
//
//	vertices = []float32{
//		1.0, -1.0,
//		0.0, 1.0,
//		1.0, -1.0,
//	}
)

const (
	FPS = 60
)

func init() {

}

func main() {
	runtime.LockOSThread()

	//TODO: Create a GL struct that handles the gl library through composition

	//GL, err := util.NewGL(util.WIDTH, util.HEIGHT, "Test Window", vertices)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer GL.Terminate()

	// Set clear to white
	//gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	// Run main render loop

}
