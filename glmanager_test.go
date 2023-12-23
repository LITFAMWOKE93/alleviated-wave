package main

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/assert"
)

func TestGLManager_SetShaderSource(t *testing.T) {
	manager := GLManager{}

	// Test setting vertex shader source
	manager.SetShaderSrouce("vertex_shader_code", "vertex")
	assert.Equal(t, "vertex_shader_code", manager.VertexShaderSource())

	// Test setting fragment shader source
	manager.SetShaderSrouce("fragment_shader_code", "fragment")
	assert.Equal(t, "fragment_shader_code", manager.FragmentShaderSource())

	// Test setting an unsupported shader type
	manager.SetShaderSrouce("unsupported_shader_code", "geometry")
	assert.Equal(t, "fragment_shader_code", manager.FragmentShaderSource())
	assert.Equal(t, "vertex_shader_code", manager.VertexShaderSource())
}

func TestGLManager_SetVertices(t *testing.T) {
	manager := GLManager{}

	// Test setting vertices
	vertices := []mgl32.Vec3{
		mgl32.Vec3{1.0, 2.0, 3.0},
		mgl32.Vec3{4.0, 5.0, 6.0},
	}

	manager.SetVertices(vertices)
	assert.Equal(t, vertices, manager.Vertices())
}

func TestGLManager_SetFloat32Vertices(t *testing.T) {
	manager := GLManager{}

	// Test setting float32 vertices after setting vertices
	vertices := []mgl32.Vec3{
		mgl32.Vec3{1.0, 2.0, 3.0},
		mgl32.Vec3{4.0, 5.0, 6.0},
	}

	manager.SetVertices(vertices)
	manager.SetFloat32Vertices()
	assert.Equal(t, []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}, manager.float32vertices)

	manager.ClearVertices()

	// Test setting float32 vertices without setting vertices
	manager.ClearFloat32Vertices()
	manager.SetFloat32Vertices()
	assert.Empty(t, manager.float32vertices)
}

func TestGLManager_BindProgram(t *testing.T) {
	manager := GLManager{}

	// Test binding program with a non-zero program value
	manager.program = 123
	manager.BindProgram()
	// Assert that the program is in use (not checking OpenGL state in this example)
}

func TestGLManager_BindVAO(t *testing.T) {
	manager := GLManager{}

	// Test binding VAO with a non-zero VBO value
	manager.vbo = 456
	manager.BindVAO()
	// Assert that the VAO is created and bound (not checking OpenGL state in this example)
}

func TestGLManager_BindVBO(t *testing.T) {
	manager := GLManager{}

	// Test binding VBO
	manager.BindVBO()
	// Assert that the VBO is created and bound (not checking OpenGL state in this example)
}

func TestGLManager_ConvertVec3ToFloat32(t *testing.T) {
	manager := GLManager{}

	// Test converting Vec3 to float32
	vertices := []mgl32.Vec3{
		mgl32.Vec3{1.0, 2.0, 3.0},
		mgl32.Vec3{4.0, 5.0, 6.0},
	}

	manager.vertices = vertices
	result := manager.ConvertVec3ToFloat32()
	assert.Equal(t, []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}, result)
}

func TestGLManager_ClearVertices(t *testing.T) {
	manager := GLManager{}

	vertices := []mgl32.Vec3{
		mgl32.Vec3{1.0, 2.0, 3.0},
		mgl32.Vec3{4.0, 5.0, 6.0},
	}

	manager.vertices = vertices
	result := manager.vertices
	assert.Equal(t, []mgl32.Vec3{mgl32.Vec3{1.0, 2.0, 3.0}, mgl32.Vec3{4.0, 5.0, 6.0}}, result)

	manager.ClearVertices()
	assert.Empty(t, manager.vertices)

}

func TestGLManager_ClearFloat32Vertices(t *testing.T) {

	manager := GLManager{}

	vertices := []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}

	manager.float32vertices = vertices
	result := manager.float32vertices
	assert.Equal(t, []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}, result)

	manager.ClearFloat32Vertices()
	assert.Empty(t, manager.float32vertices)

}
