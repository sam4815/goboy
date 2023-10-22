package gameboy

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	Gameboy  *Gameboy
	Renderer *sdl.Renderer
	Width    int32
	Height   int32
}

func GetKey(keyName string) Key {
	switch keyName {
	case "Up":
		return Up
	case "Down":
		return Down
	case "Left":
		return Left
	case "Right":
		return Right
	case "A":
		return A
	case "S":
		return B
	case "Return":
		return Start
	case "Right Shift":
		return Select
	}
	return Unsupported
}

func (display Display) Quit() {
	display.Renderer.Destroy()
}

func (display Display) HandleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			keyEvent := event.(*sdl.KeyboardEvent)

			if keyEvent.Type == sdl.KEYDOWN {
				keyName := sdl.GetKeyName(keyEvent.Keysym.Sym)
				display.Gameboy.MMU.KeyDown(GetKey(keyName))
			} else if keyEvent.Type == sdl.KEYUP {
				keyName := sdl.GetKeyName(keyEvent.Keysym.Sym)
				display.Gameboy.MMU.KeyUp(GetKey(keyName))
			}
		case *sdl.QuitEvent:
			display.Quit()
			break
		}
	}
}

func NewDisplay() *Display {
	display := Display{Height: 144, Width: 160}

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to iniialize sdl: %s\n", err)
	}

	window, _ := sdl.CreateWindow("Goboy", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		display.Width, display.Height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
	}

	display.Renderer = renderer

	return &display
}

func (display Display) Flush(canvas []byte) {
	for x := int32(0); x < display.Width; x++ {
		for y := int32(0); y < display.Height; y++ {
			pixelOffset := (x + y*display.Width) * 4
			colour := canvas[pixelOffset : pixelOffset+4]

			display.Renderer.SetDrawColor(colour[0], colour[1], colour[2], colour[3])
			display.Renderer.DrawPoint(x, y)
		}
	}
	display.Renderer.Present()
}
