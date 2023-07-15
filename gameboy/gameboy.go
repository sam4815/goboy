package gameboy

type Gameboy struct {
	Memory []byte
	CPU    CPU
}

func (gameboy Gameboy) Run() {
	// for {
	gameboy.CPU.Step(&gameboy.Memory)
	// }
}
