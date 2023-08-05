package gameboy

import "github.com/go-gl/glfw/v3.3/glfw"

type Display struct {
	Window  *glfw.Window
	Gameboy *Gameboy
}

func GetKey(key glfw.Key) Key {
	switch key {
	case glfw.KeyUp:
		return Up
	case glfw.KeyDown:
		return Down
	case glfw.KeyLeft:
		return Left
	case glfw.KeyRight:
		return Right
	case glfw.KeyZ:
		return A
	case glfw.KeyX:
		return B
	case glfw.KeyEnter:
		return Start
	case glfw.KeyRightShift:
		return Select
	}
	return Unsupported
}

func NewDisplay() Display {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(160, 144, "Gameboy", nil, nil)
	if err != nil {
		panic(err)
	}

	display := Display{Window: window}

	onKeyPress := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		gameboyKey := GetKey(key)
		if gameboyKey == Unsupported {
			return
		}

		switch action {
		case glfw.Press:
			display.Gameboy.MMU.KeyDown(gameboyKey)
		case glfw.Release:
			display.Gameboy.MMU.KeyUp(gameboyKey)
		}
	}

	window.SetKeyCallback(onKeyPress)

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}

	return display
}
