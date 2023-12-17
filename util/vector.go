package util

type Vec32 struct {
	X, Y, Z float32
}

func NewVec32(x float32, y float32, z float32) Vec32 {
	result := Vec32{X: x, Y: y, Z: z}

	return result
}

func (Vec32) Add(v1, v2 Vec32) Vec32 {
	return Vec32{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func Mult(scalar float32, v Vec32) Vec32 {
	return Vec32{scalar * v.X, scalar * v.Y, scalar * v.Z}
}
