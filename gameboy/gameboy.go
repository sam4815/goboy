package gameboy

type Gameboy struct {
	MMU MMU
	CPU CPU
}

func New(bytes []byte) Gameboy {
	gameboy := Gameboy{
		MMU: NewMMU(bytes),
		CPU: CPU{},
	}
	gameboy.CPU.Gameboy = &gameboy

	return gameboy
}

func (gameboy Gameboy) Run() {
	// for {
	gameboy.CPU.Step()
	// }
}
