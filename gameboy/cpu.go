package gameboy

type CPU struct {
	Gameboy        *Gameboy
	ProgramCounter uint16
	StackPointer   uint16
	Registers      Registers
	Clocks         Clocks
}

type Clocks struct {
	M int
	T int
}

func (cpu CPU) Immediate8() uint8 {
	return cpu.Gameboy.MMU.ReadByte(cpu.ProgramCounter + 1)
}

func (cpu CPU) Immediate16() uint16 {
	return cpu.Gameboy.MMU.ReadWord(cpu.ProgramCounter + 1)
}

func (cpu *CPU) Step() {
	opcode := DecodeOpcode(cpu.Gameboy.MMU.ReadByte(cpu.ProgramCounter))

	mTime, tTime := cpu.Execute(opcode)

	cpu.Clocks.M += mTime
	cpu.Clocks.T += tTime
}
