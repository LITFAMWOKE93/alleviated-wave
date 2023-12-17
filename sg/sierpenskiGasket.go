package main

import (
	"LITFAMWOKE93/alleviated-wave/util"
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
)

var (
	// a slice of float32, the datatype that is always fed to openGL
	// 0, 0, 0 is the center axis of the view
	SierpenskiVertices []float32
)

type Vec struct {
	X, Y, Z float32
}

func Add(v1, v2 Vec) Vec {

	return Vec{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func (v1 *Vec) Mult(scalar float32) Vec {

	return Vec{
		X: v1.X * scalar,
		Y: v1.Y * scalar,
		Z: v1.Z * scalar,
	}

}

const (
	FPS = 60

	// How many points to generate
	NUMPOSITIONS = 50
)

var (
	points []float32
)

func init() {
	// Random seed
	rand.NewSource(time.Now().Unix())

	// On init do the vector math for all of the points
	// and create an object array fom them
	v1 := Vec{-1.0, -1.0, 0.0}
	v2 := Vec{0.0, 1.0, 0.0}
	v3 := Vec{1.0, -1.0, 0.0}

	var vertices []Vec

	vertices = append(vertices, v1, v2, v3)

	var u = Add(v1, v2)
	var v = Add(v1, v3)

	var yewVee = Add(u, v)

	var p = yewVee.Mult(0.5)
	// A slice of vectors that only has the starting position
	positions := []Vec{p}
	for i := 0; i < NUMPOSITIONS-1; i++ {
		var j = math.Floor(3 * rand.Float64())
		fmt.Println(j)

		p = Add(positions[i], vertices[int(j)])
		p = p.Mult(0.5)
		positions = append(positions, p)
	}

	for _, pos := range positions {
		points = append(points, pos.X, pos.Y, pos.Z)
	}

}

func main() {
	runtime.LockOSThread()

	//Create GL instance
	GL, err := util.NewGL(800, 600, "Sierpinski Gastet", points)
	if err != nil {
		log.Fatal(err)
	}
	defer GL.Terminate()

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	if GL.VAOarray != nil {
		GL.Run(GL.RenderSierpinksiGasket)
	} else {
		fmt.Println("Array storage is empty.")
	}

	if err := gl.GetError(); err != 0 {
		log.Printf("openGL error: %v\n", err)
	}

}
