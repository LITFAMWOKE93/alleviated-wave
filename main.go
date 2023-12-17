package main

import (
	"LITFAMWOKE93/alleviated-wave/util"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
)

var (
	// a slice of float32, the datatype that is always fed to openGL
	// 0, 0, 0 is the center axis of the view
	vertices = []float32{
		1.0, -1.0,
		0.0, 1.0,
		1.0, -1.0,
	}
)

const (
	FPS = 60
)

func init() {

}

func main() {
	runtime.LockOSThread()

	GL, err := util.NewGL(util.WIDTH, util.HEIGHT, "Test Window", vertices)
	if err != nil {
		log.Fatal(err)
	}
	defer GL.Terminate()

	// Set clear to white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	// Run main render loop
	GL.Run(GL.RenderTriangle)

}
