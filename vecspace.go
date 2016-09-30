package main

import "github.com/veandco/go-sdl2/sdl"
import "github.com/veandco/go-sdl2/sdl_gfx"
import "github.com/ungerik/go3d/vec2"
import "github.com/ungerik/go3d/vec3"
import "github.com/ungerik/go3d/vec4"
import "github.com/ungerik/go3d/mat4"
import "fmt"

type Camera struct {
	Pos      vec3.T
	Target   vec3.T
	UpVector vec3.T
}

type Mesh struct {
	Name  string
	Verts []vec3.T
	Pos   vec3.T
	Rot   vec3.T
}

var meshes []Mesh

func makeCube() Mesh {
	var mesh Mesh
	mesh.Verts = append(mesh.Verts, vec3.T{-1, 1, 2})
	mesh.Verts = append(mesh.Verts, vec3.T{1, 1, 2})
	mesh.Verts = append(mesh.Verts, vec3.T{-1, -1, 2})
	mesh.Verts = append(mesh.Verts, vec3.T{-1, -1, 1})
	mesh.Verts = append(mesh.Verts, vec3.T{-1, 1, 1})
	mesh.Verts = append(mesh.Verts, vec3.T{1, 1, 1})
	mesh.Verts = append(mesh.Verts, vec3.T{1, -1, 2})
	mesh.Verts = append(mesh.Verts, vec3.T{1, -1, 1})
	mesh.Pos = vec3.T{0, 0, 0}
	mesh.Rot = vec3.T{0, 0, 0}
	mesh.Name = "Cube"
	return mesh
}

func lookAt(cam Camera) mat4.T {
	/*
		zaxis = normal(cameraTarget - cameraPosition)
		xaxis = normal(cross(cameraUpVector, zaxis))
		yaxis = cross(zaxis, xaxis)

		 xaxis.x           yaxis.x           zaxis.x          0
		 xaxis.y           yaxis.y           zaxis.y          0
		 xaxis.z           yaxis.z           zaxis.z          0
		-dot(xaxis, cameraPosition)  -dot(yaxis, cameraPosition)  -dot(zaxis, cameraPosition)  1
	*/

	zaxis := cam.Target.Sub(&cam.Pos).Normal()
	xaxis := vec3.Cross(&cam.UpVector, &zaxis)
	xaxis = xaxis.Normal()
	yaxis := vec3.Cross(&zaxis, &xaxis)

	m := mat4.T{
		vec4.T{xaxis[0], xaxis[1], xaxis[2], -vec3.Dot(&xaxis, &cam.Pos)},
		vec4.T{yaxis[0], yaxis[1], yaxis[2], -vec3.Dot(&yaxis, &cam.Pos)},
		vec4.T{zaxis[0], zaxis[1], zaxis[2], -vec3.Dot(&zaxis, &cam.Pos)},
		vec4.T{0, 0, 0, 1},
	}
	return m
}

func simpleMat() mat4.T {
	m := mat4.T{
		vec4.T{1, 0, 0, 0},
		vec4.T{0, 1, 0, 0},
		vec4.T{0, 0, 1, 0},
		vec4.T{0, 0, 0, 1},
	}
	return m
}
func simpleMat2() mat4.T {
	m := mat4.T{
		vec4.T{1, 0, 0, 0},
		vec4.T{0, 1, 0, 0},
		vec4.T{0, 0, 1, 0},
		vec4.T{0, 0, 0, 1},
	}
	return m
}

func project(p vec3.T, mat *mat4.T, pixw float32, pixh float32) vec2.T {
	point := mat.MulVec3(&p)
	x := point[0]*pixw + pixw/2.0
	y := -point[1]*pixh + pixh/2.0
	return vec2.T{x, y}
}

// w,h = aspect
func perspectiveFov(w float32, h float32, znearPlane float32, zfarPlane float32) mat4.T {
	m := mat4.T{
		vec4.T{w, 0, 0, 0},
		vec4.T{0, h, 0, 0},
		vec4.T{0, 0, zfarPlane / (znearPlane - zfarPlane), znearPlane * zfarPlane / (znearPlane - zfarPlane)},
		vec4.T{0, 0, -1, 0},
	}

	return m
}

func main() {

	meshes = []Mesh{makeCube()}
	/*cam := Camera{
		Pos:      vec3.T{0, 0, -5.0},
		Target:   vec3.T{0, 0, 5.0},
		UpVector: vec3.T{0, 1, 0},
	}*/
	viewMatrix := simpleMat() //lookAt(cam)
	projectionMatrix := perspectiveFov(1, 1, 4, 10)
	transformMatrix := viewMatrix.MultMatrix(&projectionMatrix)

	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("vectorspace", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	/*surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}*/

	tm := simpleMat2()

	fmt.Println("running")

	var renderer *sdl.Renderer

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE); err != nil {
		fmt.Println(err)
		return
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.MouseMotionEvent:
				//fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
				//	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
				tm[3][0] = float32(t.X-400) / 100.0
				tm[3][1] = float32(t.Y-300) / 100.0

			case *sdl.KeyDownEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)

				if t.Keysym.Sym == 's' {
					tm[3][2] -= 0.1
				} else if t.Keysym.Sym == 'w' {
					tm[3][2] += 0.1
				}
			}
		}

		//rect := sdl.Rect{0, 0, 800, 600}
		//surface.FillRect(&rect, 0)

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		for _, m := range meshes {
			var lastp vec2.T

			for i, v := range m.Verts {
				v = tm.MulVec3(&v)
				p := project(v, transformMatrix, 800, 600)
				//fmt.Println(p)

				//rect := sdl.Rect{int32(p[0]), int32(p[1]), 1, 1}
				//surface.FillRect(&rect, 0xffffffff)

				if i > 0 {
					gfx.LineRGBA(renderer, int(lastp[0]), int(lastp[1]), int(p[0]), int(p[1]), 0xff, 0xff, 0xff, 0xff)
				}
				lastp = p
			}
		}

		renderer.Present()
	}

	sdl.Delay(2000)
	sdl.Quit()
}
