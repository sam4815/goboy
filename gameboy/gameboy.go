package gameboy

type Gameboy struct {
	MMU      MMU
	CPU      CPU
	Metadata Metadata
}

func New(bytes []byte) Gameboy {
	gameboy := Gameboy{
		MMU:      NewMMU(bytes),
		CPU:      CPU{},
		Metadata: ParseMetadata(bytes),
	}
	gameboy.CPU.Gameboy = &gameboy

	return gameboy
}

func (gameboy Gameboy) Run() {
	// for {
	gameboy.CPU.Step()
	// }
}
