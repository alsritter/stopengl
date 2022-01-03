package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
	"stopengl/gfx"
)

const (
	WIDTH  = 500
	HEIGHT = 500
)

var (
	vertices = []float32{
		// 位置              // 颜色
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // 右下
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // 左下
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // 顶部
	}

	indices = []uint32{ // 注意索引从0开始!
		0, 1, 2, // 第一个三角形
	}
)

func main() {
	runtime.LockOSThread()
	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()
	defer program.Delete()

	// 必须告诉 OpenGL 渲染窗口的尺寸大小，即视口(Viewport)，这样 OpenGL 才只能知道怎样根据窗口大小显示数据和坐标。
	gl.Viewport(0, 0, WIDTH, HEIGHT) // 起点为左下角
	// 窗口大小被改变的回调函数
	window.SetFramebufferSizeCallback(framebuffer_size_callback)

	vao, vbo, ebo := makeVao(vertices, indices)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL) // LINE 是线框模式

	//渲染循环
	for !window.ShouldClose() {
		//用户输入
		processInput(window)
		glfw.PollEvents()
		draw(vao, window, program)
	}

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}

func draw(vao uint32, window *glfw.Window, program *gfx.Program) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0) //状态设置
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// draw loop
	program.Use()
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil) // 第二个参数是count
	gl.BindVertexArray(0)
	// end of draw loop

	window.SwapBuffers()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("glfw initializes err: %v", err))
	}
	// 通过这些枚举值来设置 GLFW 的参数
	glfw.WindowHint(glfw.Resizable, glfw.False)                 // 设置窗口大小无法修改
	glfw.WindowHint(glfw.ContextVersionMajor, 3)                // OpenGL最大版本
	glfw.WindowHint(glfw.ContextVersionMinor, 3)                // OpenGl 最小版本
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) // 明确核心模式
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)    // Mac使用时需要加上
	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "LearnOpenGL", nil, nil)
	if window == nil || err != nil {
		panic(err)
	}
	log.Println("created window")
	window.MakeContextCurrent() // 通知 glfw 将当前窗口上下文绑定到当前线程的上下文
	return window
}

func printOpenGLInfo() {
	log.Println("=========================")
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	var nrAttributes int32
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nrAttributes)
	log.Printf("Number of Vertices currently supported: %d", nrAttributes)
	log.Println("=========================")
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() *gfx.Program {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	printOpenGLInfo()
	vertShader, err := gfx.NewShaderFromFile("shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := gfx.NewShaderFromFile("shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shaderProgram, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}
	return shaderProgram
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32, idxs []uint32) (uint32, uint32, uint32) {
	var vbo, ebo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	gl.GenVertexArrays(1, &vao)

	// 1. 绑定顶点数组对象
	gl.BindVertexArray(vao)

	// 2. 把我们的顶点数组复制到一个顶点缓冲中，供 OpenGL 使用
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	// 3. 复制我们的索引数组到一个索引缓冲中，供 OpenGL 使用
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(idxs), gl.Ptr(idxs), gl.STATIC_DRAW)

	// 4. 设定顶点属性指针
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0) // VertexAttribPointer 偏移已经被废弃
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
	gl.EnableVertexAttribArray(1)
	return vao, vbo, ebo
}

// 监听进程输入
func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		log.Println("escape pressed")
		window.SetShouldClose(true)
	}
}

// 在调整指定窗口的帧缓冲区大小时调用。
func framebuffer_size_callback(window *glfw.Window, width int, height int) {
	log.Printf("resize width:%d, height:%d", width, height)
	gl.Viewport(0, 0, int32(width), int32(height))
}
