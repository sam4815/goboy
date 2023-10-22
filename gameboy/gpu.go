package gameboy

import "log"

type GPUMode byte

const (
	HorizontalBlank GPUMode = 0
	VerticalBlank   GPUMode = 1
	OAMScanline     GPUMode = 2
	VRAMScanline    GPUMode = 3
)

type GPU struct {
	Gameboy           *Gameboy
	Mode              GPUMode
	Clock             int
	Line              byte
	VRAM              []byte
	Canvas            []byte
	ScrollX           byte
	ScrollY           byte
	TileSet           [][][]byte
	BackgroundMap     bool
	BackgroundTile    uint8
	BackgroundPalette [][]byte
	SwitchLCD         bool
	SwitchBackground  bool
}

func NewGPU() *GPU {
	gpu := GPU{
		Canvas:            make([]byte, 160*144*4),
		VRAM:              make([]byte, 8192),
		BackgroundPalette: make([][]byte, 4),
		TileSet:           make([][][]byte, 512),
	}

	for i := 0; i < 512; i++ {
		gpu.TileSet[i] = make([][]byte, 8)
		for j := 0; j < 8; j++ {
			gpu.TileSet[i][j] = []byte{0, 0, 0, 0, 0, 0, 0, 0}
		}
	}

	for i := 0; i < 160*144*4; i++ {
		gpu.Canvas[i] = 255
	}

	return &gpu
}

func (gpu *GPU) Sum() int {
	sum := 0
	for _, num := range gpu.Canvas {
		sum += int(num)
	}
	return sum
}

func (gpu *GPU) Step() {
	switch gpu.Mode {
	case OAMScanline:
		if gpu.Clock >= 80 {
			gpu.Clock = 0
			gpu.Mode = VRAMScanline
		}

	case VRAMScanline:
		if gpu.Clock >= 172 {
			gpu.Clock = 0
			gpu.Mode = HorizontalBlank

			gpu.WriteScanline()
		}

	case HorizontalBlank:
		if gpu.Clock >= 204 {
			gpu.Clock = 0
			gpu.Line += 1

			if gpu.Line == 143 {
				log.Print("FLUSHING: ", gpu.Sum())
				gpu.Mode = 1
				gpu.Gameboy.Display.Flush(gpu.Canvas)
			} else {
				gpu.Mode = 2
			}
		}

	case VerticalBlank:
		if gpu.Clock >= 456 {
			gpu.Clock = 0
			gpu.Line += 1

			if gpu.Line > 153 {
				gpu.Line = 0
				gpu.Mode = 2
			}
		}
	}
}

func (gpu *GPU) UpdateTile(address uint16, value byte) {
	tile := (address >> 4) & 511
	y := (address >> 1) & 7

	for x := 0; x < 8; x++ {
		tileValue := ((gpu.VRAM[address] >> (7 - x)) & 1)
		tileValue += ((gpu.VRAM[address] >> (7 - x)) & 1) * 2

		gpu.TileSet[tile][y][x] = tileValue
	}
}

func (gpu *GPU) WriteScanline() {
	mapOffset := uint16(0x1800)

	if gpu.BackgroundMap {
		mapOffset = 0x1C00
	}

	currentLine := gpu.Line + gpu.ScrollY
	mapOffset += uint16((currentLine & 255) >> 3)

	lineOffset := gpu.ScrollX >> 3

	y := currentLine & 7
	x := gpu.ScrollX & 7

	canvasOffset := uint32(gpu.Line) * 160 * 4

	tile := uint16(gpu.VRAM[mapOffset+uint16(lineOffset)])
	if gpu.BackgroundTile == 1 && tile < 128 {
		tile += 256
	}

	for i := 0; i < 160; i++ {
		tilePixel := gpu.TileSet[tile][y][x]
		colour := gpu.BackgroundPalette[tilePixel]

		copy(gpu.Canvas[canvasOffset:], colour)
		canvasOffset += 4

		x += 1
		if x == 8 {
			x = 0
			lineOffset += 1

			tile := uint16(gpu.VRAM[mapOffset+uint16(lineOffset)])
			if gpu.BackgroundTile == 1 && tile < 128 {
				tile += 256
			}
		}
	}
}

func (gpu GPU) LCDByte() byte {
	lcdByte := byte(0)

	if gpu.SwitchLCD {
		lcdByte |= (1 << 7)
	}
	if gpu.BackgroundTile == 1 {
		lcdByte |= (1 << 4)
	}
	if gpu.BackgroundMap {
		lcdByte |= (1 << 3)
	}
	if gpu.SwitchBackground {
		lcdByte |= 1
	}

	return lcdByte
}

func (gpu *GPU) SetLCDByte(b byte) {
	gpu.SwitchLCD = (b>>7)&1 == 1
	gpu.BackgroundTile = (b >> 4) & 1
	gpu.BackgroundMap = (b>>3)&1 == 1
	gpu.SwitchBackground = b&1 == 1
}

func (gpu GPU) ReadRegisters(address uint16) byte {
	switch address {
	case 0xFF40:
		return gpu.LCDByte()
	case 0xFF42:
		return gpu.ScrollY
	case 0xFF43:
		return gpu.ScrollX
	case 0xFF44:
		return gpu.Line
	}
	return 0
}

func (gpu *GPU) WriteRegisters(address uint16, value byte) {
	switch address {
	case 0xFF40:
		gpu.SetLCDByte(value)
	case 0xFF42:
		gpu.ScrollY = value
	case 0xFF43:
		gpu.ScrollX = value
	case 0xFF47:
		for i := 0; i < 4; i++ {
			// log.Print("hm", value)
			switch (value >> (i * 2)) & 0b11 {
			case 0:
				gpu.BackgroundPalette[i] = []byte{255, 255, 255, 255}
			case 1:
				gpu.BackgroundPalette[i] = []byte{192, 192, 192, 255}
			case 2:
				gpu.BackgroundPalette[i] = []byte{96, 96, 96, 255}
			case 3:
				gpu.BackgroundPalette[i] = []byte{0, 0, 0, 255}
			}
		}
	}
}
