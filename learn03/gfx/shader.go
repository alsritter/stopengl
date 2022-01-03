package gfx

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"io/ioutil"
	"strings"
)

type Shader struct {
	handle uint32
}

type Program struct {
	handle  uint32
	shaders []*Shader
}

func (shader *Shader) Delete() {
	gl.DeleteShader(shader.handle)
}

func (prog *Program) Delete() {
	for _, shader := range prog.shaders {
		shader.Delete()
	}
	gl.DeleteProgram(prog.handle)
}

func (prog *Program) Attach(shaders ...*Shader) {
	for _, shader := range shaders {
		gl.AttachShader(prog.handle, shader.handle)
		prog.shaders = append(prog.shaders, shader)
	}
}

func (prog *Program) SetUniformF4(name string, v1, v2, v3, v4 float32) {
	location := gl.GetUniformLocation(prog.handle, gl.Str(name+"\x00"))
	gl.Uniform4f(location, v1, v2, v3, v4)
}

func (prog *Program) Use() {
	gl.UseProgram(prog.handle)
}

func (prog *Program) Link() error {
	gl.LinkProgram(prog.handle)
	return getGlError(prog.handle, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog,
		"PROGRAM::LINKING_FAILURE")
}

func NewProgram(shaders ...*Shader) (*Program, error) {
	prog := &Program{handle: gl.CreateProgram()}
	prog.Attach(shaders...)
	if err := prog.Link(); err != nil {
		return nil, err
	}
	return prog, nil
}

func NewShader(src string, sType uint32) (*Shader, error) {
	handle := gl.CreateShader(sType)
	glSrc, freeFn := gl.Strs(src + "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrc, nil)
	gl.CompileShader(handle)
	err := getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"SHADER::COMPILE_FAILURE::")
	if err != nil {
		return nil, err
	}
	return &Shader{handle: handle}, nil
}

func NewShaderFromFile(file string, sType uint32) (*Shader, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	handle := gl.CreateShader(sType)
	glSrc, freeFn := gl.Strs(string(src) + "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrc, nil)
	gl.CompileShader(handle)
	err = getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"SHADER::COMPILE_FAILURE::"+file)
	if err != nil {
		return nil, err
	}
	return &Shader{handle: handle}, nil
}

// 取得当前 handle 对象的入口函数
type getObjIv func(uint32, uint32, *int32)

// 取得 Log 的入口函数
type getObjInfoLog func(uint32, int32, *int32, *uint8)

//  getGlError
//  @Description: 取得错误消息
//  @param glHandle 当前操作对象的 ID
//  @param checkTrueParam 当前正在执行的操作
//  @param getObjIvFn 取得当前 handle 对象的入口函数
//  @param getObjInfoLogFn 取得 Log 的入口函数
//  @param failMsg 错误消息
//  @return error
//
func getGlError(
	glHandle uint32,
	checkTrueParam uint32,
	getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog,
	failMsg string) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, log)

		return fmt.Errorf("%s: %s", failMsg, gl.GoStr(log))
	}

	return nil
}
