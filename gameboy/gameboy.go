package gameboy

import "log"

type Gameboy struct {
	MMU      MMU
	CPU      CPU
	Metadata Metadata
}

func New(bytes []byte) Gameboy {
	gameboy := Gameboy{
		MMU:      NewMMU(bytes),
		CPU:      NewCPU(),
		Metadata: ParseMetadata(bytes),
	}
	gameboy.CPU.Gameboy = &gameboy

	return gameboy
}

func (gameboy *Gameboy) Run() {
	for i := 0; i < 10; i++ {
		gameboy.CPU.Step()
		log.Print(gameboy.CPU.ProgramCounter, gameboy.CPU.StackPointer, gameboy.CPU.Registers)
	}
}
