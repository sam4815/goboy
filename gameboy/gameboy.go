package gameboy

import (
	"encoding/hex"
	"fmt"
)

type Gameboy struct {
	MMU      MMU
	CPU      CPU
	Metadata Metadata
	Display  Display
}

func New(bytes []byte) Gameboy {
	gameboy := Gameboy{
		MMU:      NewMMU(bytes),
		CPU:      NewCPU(),
		Metadata: ParseMetadata(bytes),
		Display:  NewDisplay(),
	}
	gameboy.CPU.Gameboy = &gameboy
	gameboy.Display.Gameboy = &gameboy

	return gameboy
}

func (gameboy *Gameboy) Run() {
	for {
		gameboy.CPU.Step()
		if gameboy.CPU.ProgramCounter > 355 {
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
	}
}
