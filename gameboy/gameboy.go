package gameboy

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
	gameboy.CPU.ProgramCounter = 0x150

	for i := 0; i < 100; i++ {
		gameboy.CPU.Step()
	}
}
