package main

import (
	"fmt"
	"strings"
	"time"
	"math"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/go-gl/gl/v4.6-core/gl"
	"math/rand"
)


var hexSize float32

func main() {
	tau := math.Pi * 2

	xhexConsts := [6]float32{1,float32(math.Cos(tau/6)),float32(math.Cos(tau/3)),-1,float32(math.Cos(2*tau/3)),float32(math.Cos(5*tau/6))}
	yhexConsts := [6]float32{0,float32(math.Sin(tau/6)),float32(math.Sin(tau/3)),0, float32(math.Sin(2*tau/3)),float32(math.Sin(5*tau/6))}
	var hexConsts [12]float32
	for i:=0;i<6;i++ {
		hexConsts[i*2] = xhexConsts[i]
		hexConsts[i*2+1] = yhexConsts[i]
	}
	
	
	hexSize = 1
	
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {panic(err)}

	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", 200,200,400,400,sdl.WINDOW_OPENGL)

	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println(version)

	
	vertexShaderSource := `#version 460 core
	layout (location = 0) in vec2 pos;

	out vec2 col;

	void main() {	
		gl_Position = vec4(pos, 0.0, 1.0);
		col = vec2(pos)/2 + 0.5;
	}`+"\x00"

	fragShaderSource := `#version 460 core
	out vec4 fragColor;
	in vec2 col;

	void main() {
		fragColor = vec4(col.x, 1.0 - (col.x + col.y), col.y, 1.0);
	}`+"\x00"

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)

	csrc,free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, csrc, nil)
	free()
	gl.CompileShader(vertexShader)
	var status int32
	gl.GetShaderiv(vertexShader,gl.COMPILE_STATUS,&status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader,gl.INFO_LOG_LENGTH,&logLength)
		log:=strings.Repeat("\x00",int(logLength+1))
		gl.GetShaderInfoLog(vertexShader,logLength,nil,gl.Str(log))
		panic("vertexshader\n"+log)
	}

	
	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)

	csrc2,free := gl.Strs(fragShaderSource)
	
	gl.ShaderSource(fragShader, 1, csrc2, nil)
	free()
	gl.CompileShader(fragShader)
	var status2 int32
	gl.GetShaderiv(fragShader,gl.COMPILE_STATUS,&status2)
	if status2 == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragShader,gl.INFO_LOG_LENGTH,&logLength)
		log:=strings.Repeat("\x00",int(logLength+1))
		gl.GetShaderInfoLog(fragShader,logLength,nil,gl.Str(log))
		panic("fragshader\n" + log)
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragShader)
	gl.LinkProgram(shaderProgram)

	var status3 int32
	gl.GetProgramiv(shaderProgram,gl.LINK_STATUS,&status3)
	if status3 == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram,gl.INFO_LOG_LENGTH,&logLength)
		log:=strings.Repeat("\x00",int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram,logLength,nil,gl.Str(log))
		panic("shader3\n" + log)
	}




	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragShader)

	

	//updateHexRadius(0.5, &hexConsts)
	
	
	
	


	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		gl.ClearColor(0,0,0,0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		updateHexRadius(rand.Float32(), &hexConsts)
		drawHex(rand.Float32()*2-1,rand.Float32()*2-1,&hexConsts, shaderProgram)
		window.GLSwap()

		//500ms
		time.Sleep(500000000)
		//vertices[0] += 0.1
		//fmt.Println(vertices[0]);
	}

}

func updateHexRadius(newSize float32, b *[12]float32) {
	a := newSize/hexSize
	hexSize = newSize
	for i := 0; i < 12; i++ {
		(*b)[i] *= a
	}

	return
}

//updateHexRadius(10)

func drawHex(x float32,y float32, a *[12]float32, b uint32) {
	var vertices [12]float32
	for i:= 0; i < 12; i+=2 {
		vertices[i] = x + (*a)[i]
		vertices[i+1] = y + (*a)[i+1]
	}

	indices := []uint32{
		0,1,2,
		0,2,3,
		0,3,4,
		0,4,5}
	
	var VBO uint32
	var VAO uint32
	var EBO uint32

	gl.GenBuffers(1,&VBO)
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &EBO)
	
	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,gl.Ptr(&vertices[0]),gl.STATIC_DRAW)
	
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(&indices[0]), gl.STATIC_DRAW)
	
	gl.VertexAttribPointer(0,2,gl.FLOAT,false,2*4, nil)
	gl.EnableVertexAttribArray(0)
	
	gl.UseProgram(b)
	gl.BindVertexArray(VAO)
	//gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.DrawElements(gl.TRIANGLES, 12, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
	
	return
}