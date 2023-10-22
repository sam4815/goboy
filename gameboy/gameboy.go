package gameboy

import (
	"encoding/hex"
	"fmt"
)

type Gameboy struct {
	MMU      *MMU
	CPU      *CPU
	GPU      *GPU
	Display  *Display
	Metadata Metadata
}

func New(bytes []byte) Gameboy {
	gameboy := Gameboy{
		MMU:      NewMMU(bytes),
		CPU:      NewCPU(),
		GPU:      NewGPU(),
		Metadata: ParseMetadata(bytes),
		Display:  NewDisplay(),
	}
	gameboy.CPU.Gameboy = &gameboy
	gameboy.GPU.Gameboy = &gameboy
	gameboy.MMU.Gameboy = &gameboy
	gameboy.Display.Gameboy = &gameboy

	return gameboy
}

func (gameboy *Gameboy) Run() {
	for i := 0; i < 1000000; i++ {
		gameboy.Display.HandleEvents()
		gameboy.StepFrame()
	}
	// for {
	// 	graphix := make([]byte, 160*144*4)
	// 	for x := 0; x < 160; x++ {
	// 		for y := 0; y < 144; y++ {
	// 			idx := (x + (y * 160)) * 4
	// 			graphix[idx] = 92
	// 			graphix[idx+1] = 192
	// 			graphix[idx+2] = 92
	// 			graphix[idx+3] = 255
	// 		}
	// 	}

	// 	gameboy.Display.HandleEvents()
	// 	gameboy.Display.Flush(graphix)

	// 	time.Sleep(time.Millisecond * 200)
	// }
}

func (gameboy *Gameboy) StepFrame() {
	numCyclesInFrame := 70224
	targetCycles := gameboy.CPU.Cycles + numCyclesInFrame

	// steps := 0
	// log.Print(gameboy.CPU.ProgramCounter, gameboy.CPU.CurrentByte())
	// log.Print(gameboy.CPU.Decoder.DecodeUnprefixed(gameboy.CPU.CurrentByte()))
	for gameboy.CPU.Cycles < targetCycles {
		// log.Print(gameboy.CPU.Steps, gameboy.CPU.ProgramCounter)
		gameboy.CPU.Step()
		gameboy.GPU.Step()
		gameboy.CPU.Steps += 1
		// if gameboy.CPU.Steps >= 54872 && gameboy.CPU.Steps <= 54982 {
		// 	gameboy.Dump()
		// }
	}
}

func (gameboy *Gameboy) Dump() {
	fmt.Print(
		"OP: 0x", hex.EncodeToString([]byte{gameboy.MMU.ReadByte(gameboy.CPU.ProgramCounter)}),
		" PC: ", gameboy.CPU.ProgramCounter,
		" SP: ", gameboy.CPU.StackPointer,
		" A: ", gameboy.CPU.Registers.A,
		" B: ", gameboy.CPU.Registers.B,
		" C: ", gameboy.CPU.Registers.C,
		" D: ", gameboy.CPU.Registers.D,
		" E: ", gameboy.CPU.Registers.E,
		" H: ", gameboy.CPU.Registers.H,
		" L: ", gameboy.CPU.Registers.L,
		" Z: ", gameboy.CPU.Registers.F.Zero,
		" N: ", gameboy.CPU.Registers.F.Subtract,
		" H: ", gameboy.CPU.Registers.F.HalfCarry,
		" C: ", gameboy.CPU.Registers.F.Carry, "\n",
	)
}
