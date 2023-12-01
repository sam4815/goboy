package gameboy

import (
	"log"
)

type CPU struct {
	Gameboy            *Gameboy
	ProgramCounter     uint16
	NextProgramCounter uint16
	StackPointer       uint16
	Registers          Registers
	Cycles             int
	Steps              int
	Decoder            Decoder
	EnableInterrupts   bool
}

func NewCPU() *CPU {
	return &CPU{
		ProgramCounter:     0x150,
		NextProgramCounter: 0x01,
		StackPointer:       0xfffe,
		Cycles:             20,
		Decoder:            NewDecoder(),
		Registers: Registers{
			A: 0x01,
			F: Flags{Zero: true, Carry: true},
			B: 0x0,
			C: 0x13,
			D: 0x0,
			E: 0xD8,
			H: 0x01,
			L: 0x4D,
		},
	}
}

func (cpu CPU) CurrentByte() uint8 {
	return cpu.Gameboy.MMU.ReadByte(cpu.ProgramCounter)
}

func (cpu CPU) ImmediateByte() uint8 {
	return cpu.Gameboy.MMU.ReadByte(cpu.ProgramCounter + 1)
}

func (cpu CPU) ImmediateWord() uint16 {
	return cpu.Gameboy.MMU.ReadWord(cpu.ProgramCounter + 1)
}

func (cpu *CPU) PopStack() uint16 {
	cpu.StackPointer += 2
	return cpu.Gameboy.MMU.ReadWord(cpu.StackPointer - 2)
}

func (cpu *CPU) PushStack(value uint16) {
	cpu.StackPointer -= 2
	cpu.Gameboy.MMU.WriteWord(cpu.StackPointer, value)
}

func (cpu *CPU) GetOperand(operand OperandInfo) uint16 {
	var value uint16

	switch operand.Name {
	case "AF":
		value = cpu.Registers.GetPair(AF)
	case "BC":
		value = cpu.Registers.GetPair(BC)
	case "DE":
		value = cpu.Registers.GetPair(DE)
	case "HL":
		value = cpu.Registers.GetPair(HL)
	case "A":
		value = uint16(cpu.Registers.A)
	case "B":
		value = uint16(cpu.Registers.B)
	case "C":
		value = uint16(cpu.Registers.C)
	case "D":
		value = uint16(cpu.Registers.D)
	case "E":
		value = uint16(cpu.Registers.E)
	case "H":
		value = uint16(cpu.Registers.H)
	case "L":
		value = uint16(cpu.Registers.L)
	case "SP":
		value = cpu.StackPointer
	case "n8", "e8":
		value = uint16(cpu.ImmediateByte())
	case "a8":
		address := 0xFF00 + uint16(cpu.ImmediateByte())
		value = uint16(cpu.Gameboy.MMU.ReadByte(address))
	case "c8":
		value = 0xFF00 + uint16(cpu.Registers.C)
	case "n16", "a16":
		value = cpu.ImmediateWord()
	default:
		log.Fatal("unsupported operand named ", operand.Name)
	}

	if operand.Increment && !operand.Immediate {
		cpu.SetOperand(OperandInfo{Name: operand.Name, Immediate: true}, value+1)
	}
	if operand.Increment && operand.Immediate {
		cpu.SetOperand(OperandInfo{Name: operand.Name, Immediate: true}, uint16(cpu.ImmediateByte())+1)
	}
	if operand.Decrement && !operand.Immediate {
		cpu.SetOperand(OperandInfo{Name: operand.Name, Immediate: true}, value-1)
	}
	if !operand.Immediate {
		value = cpu.Gameboy.MMU.ReadWord(value)
	}

	return value
}

func (cpu *CPU) SetOperand(operand OperandInfo, value uint16) {
	if !operand.Immediate {
		address := cpu.GetOperand(OperandInfo{Name: operand.Name, Immediate: true})
		cpu.Gameboy.MMU.WriteWord(address, value)

		if operand.Decrement {
			cpu.SetOperand(OperandInfo{Name: operand.Name, Immediate: true}, address-1)
		}

		if operand.Increment {
			cpu.SetOperand(OperandInfo{Name: operand.Name, Immediate: true}, address+1)
		}

		return
	}

	switch operand.Name {
	case "AF":
		cpu.Registers.SetPair(AF, value)
	case "BC":
		cpu.Registers.SetPair(BC, value)
	case "DE":
		cpu.Registers.SetPair(DE, value)
	case "HL":
		cpu.Registers.SetPair(HL, value)
	case "A":
		cpu.Registers.A = byte(value)
	case "B":
		cpu.Registers.B = byte(value)
	case "C":
		cpu.Registers.C = byte(value)
	case "D":
		cpu.Registers.D = byte(value)
	case "E":
		cpu.Registers.E = byte(value)
	case "H":
		cpu.Registers.H = byte(value)
	case "L":
		cpu.Registers.L = byte(value)
	case "SP":
		cpu.StackPointer = value
	case "a8":
		address := 0xFF00 + uint16(cpu.ImmediateByte())
		cpu.Gameboy.MMU.WriteByte(address, byte(value))
	}
}

func (cpu *CPU) Step() {
	opcode := cpu.Decoder.DecodeUnprefixed(cpu.CurrentByte())

	if opcode.Mnemonic == "PREFIX" {
		opcode = cpu.Decoder.DecodePrefixed(cpu.ImmediateByte())
	}

	cycles := cpu.Execute(opcode)

	cpu.Cycles += int(cycles)
	cpu.Gameboy.GPU.Clock += int(cycles)

	if cpu.NextProgramCounter != 0x01 {
		cpu.ProgramCounter = cpu.NextProgramCounter
		cpu.NextProgramCounter = 0x01
	} else {
		cpu.ProgramCounter += opcode.Bytes
	}
}
