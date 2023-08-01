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

func (cpu *CPU) Immediate8() uint8 {
	immediate8 := cpu.Gameboy.MMU.ReadByte(cpu.ProgramCounter)
	cpu.ProgramCounter += 1

	return immediate8
}

func (cpu *CPU) Immediate16() uint16 {
	immediate16 := cpu.Gameboy.MMU.ReadWord(cpu.ProgramCounter)
	cpu.ProgramCounter += 2

	return immediate16
}

func (cpu *CPU) Step() {
	opcode := DecodeOpcode(cpu.Gameboy.MMU.ReadByte(cpu.ProgramCounter))
	cpu.ProgramCounter += 1

	mTime, tTime := cpu.Execute(opcode)

	cpu.Clocks.M += mTime
	cpu.Clocks.T += tTime
}
