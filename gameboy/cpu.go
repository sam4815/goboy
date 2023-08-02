package gameboy

type CPU struct {
	Gameboy        *Gameboy
	ProgramCounter uint16
	StackPointer   uint16
	Registers      Registers
	Cycles         int
	Decoder        Decoder
}

func NewCPU() CPU {
	return CPU{Decoder: NewDecoder()}
}

func (cpu *CPU) CurrentByte() uint8 {
	return cpu.Gameboy.MMU.ReadByte(cpu.ProgramCounter)
}

func (cpu *CPU) ImmediateByte() uint8 {
	return cpu.Gameboy.MMU.ReadByte(cpu.ProgramCounter + 1)
}

func (cpu *CPU) ImmediateWord() uint16 {
	return cpu.Gameboy.MMU.ReadWord(cpu.ProgramCounter + 1)
}

func (cpu *CPU) Step() {
	opcode := cpu.Decoder.DecodeUnprefixed(cpu.CurrentByte())

	if opcode.Mnemonic == "PREFIX" {
		opcode = cpu.Decoder.DecodePrefixed(cpu.ImmediateByte())
	}

	cycleIndex := cpu.Execute(opcode)

	cpu.ProgramCounter += opcode.Bytes
	cpu.Cycles += int(opcode.Cycles[cycleIndex])
}
