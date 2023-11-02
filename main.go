package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type ObjKey string

const (
	Vector      ObjKey = "v"
	Normal             = "vn"
	TexCoord           = "vt"
	Face               = "f"
	ObjectName         = "o"
	UseMaterial        = "mtllib"
	MaterialLib        = "usemtl"
)

const (
	width  = 640
	height = 480
)

type Triangle struct {
	points [3]vec3
	norms  [3]vec3
	// texCoords [3]vec2
}

type Mesh struct {
	triangles []Triangle
	// textures []Texture
}

type Texture struct {
	// TODO: PBR textures
	// disney stuff
}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

type v3 struct {
	x float32
	y float32
	z float32
}

func main() {
	// get args to cmdline
	if len(os.Args) == 2 {
		parseObj(os.Args[1])
	} else {
		fmt.Println("Incorrect Argument Number. Example usage: go run . <filename>")
	}

	// here init pixel buffer
	// pixels =

	// GLFW stuff
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(width, height, "GO-RENDER", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	frameBuffer := new([width * height]v3)
	data := make([]float32, width*height*3)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			frameBuffer[i*width+j] = v3{0.5, float32(i) / height, float32(j) / width}
			data[i*3*width+3*j+0] = 1
			data[i*3*width+3*j+1] = float32(i) / height
			data[i*3*width+3*j+2] = float32(j) / width
			// frameBuffer[i*width+j] = vec3{0.5, 0.5, 0.5}
		}
	}
	// fmt.Println(data)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// define fragment shader
	fragShaderSrc := `
		#version 120

		uniform sampler2D input_tex;
		uniform vec4 BufInfo;

		void main() {
			gl_FragColor = texture2D(input_tex, gl_FragCoord.st * BufInfo.zw);
			// gl_FragColor = vec4(0.5, 0.5, 0.5, 0.4);
		}
	` + "\x00"
	cFragShaderSrc, free := gl.Strs(fragShaderSrc)
	// create shader and compile
	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragShader, 1, cFragShaderSrc, nil)
	free()
	gl.CompileShader(fragShader)

	var infoLogLen int32
	gl.GetShaderiv(fragShader, gl.INFO_LOG_LENGTH, &infoLogLen)
	var infoLog = make([]byte, infoLogLen+1)
	gl.GetShaderInfoLog(fragShader, infoLogLen, nil, &infoLog[0])
	fmt.Println(string(infoLog))

	// combine shaders into program
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, fragShader)
	gl.LinkProgram(shaderProgram)

	// create texture
	var texture uint32
	gl.ActiveTexture(gl.TEXTURE0)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// set texture params
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP)
	// allocate texImage
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB32F_ARB, width, height, 0, gl.LUMINANCE, gl.FLOAT, gl.Ptr(nil))
	// ??? >> Initialize some gl state
	gl.Disable(gl.DEPTH_TEST)

	gl.UseProgram(shaderProgram)

	// set the inputs
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram, gl.Str("input_tex\x00")), 0)

	dims := new([4]int32)
	gl.GetIntegerv(gl.VIEWPORT, &dims[0])
	bufInfo := [4]float32{float32(dims[2]), float32(dims[3]), 1.0 / float32(dims[2]), 1.0 / float32(dims[3])}
	gl.Uniform4fv(gl.GetUniformLocation(shaderProgram, gl.Str("BufInfo\x00")), 1, &bufInfo[0])

	// IDK what these are for
	// gl.ActiveTexture(gl.TEXTURE0)
	// gl.BindTexture(gl.TEXTURE_2D, texture)

	// gl.MatrixMode(gl.PROJECTION)
	// gl.LoadIdentity()

	// gl.MatrixMode(gl.MODELVIEW)
	// gl.LoadIdentity()

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		glfw.PollEvents()
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0,
			width, height, gl.RGB, gl.FLOAT, gl.Ptr(data))
		gl.Recti(1, 1, -1, -1)
		window.SwapBuffers()
	}
}

func parseObj(filepath string) (int, error) {
	// check if file ends in .obj
	if !strings.HasSuffix(filepath, ".obj") {
		fmt.Println("Incorrect file type. Provide an .obj file")
		return 1, nil
	}
	// try to open file
	file, err := os.Open(filepath)
	if err != nil {
		return 1, err
	}
	defer file.Close()

	// save some info
	// vertices := []vec3{}
	// normals := []vec3{}
	// texCoords := []vec2{}
	// ultimately there should be a list of meshes returned

	// read through each line of the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v := strings.Fields(scanner.Text())

		// cases based on what kind of keyword
		switch ObjKey(v[0]) {
		case Vector:
		case Normal:
		case TexCoord:
		case Face:
		case ObjectName:
		case UseMaterial:
		case MaterialLib:
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return 1, err
	}

	return 0, nil
}
