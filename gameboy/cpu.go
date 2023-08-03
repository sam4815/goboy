package gameboy

import (
	"log"
)

type CPU struct {
	Gameboy        *Gameboy
	ProgramCounter uint16
	StackPointer   uint16
	Registers      Registers
	Cycles         int
	Decoder        Decoder
}

func NewCPU() CPU {
	return CPU{
		ProgramCounter: 0x150,
		StackPointer:   0xfffe,
		Decoder:        NewDecoder(),
		Registers: Registers{
			A: 0x11,
			F: Flags{Zero: true},
			B: 0x0,
			C: 0x0,
			D: 0xff,
			E: 0x56,
			H: 0x0,
			L: 0xd,
		},
	}
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

func (cpu CPU) GetOperand(operand OperandInfo) uint16 {
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
	case "n8":
		value = uint16(cpu.ImmediateByte())
	case "n16":
		value = cpu.ImmediateWord()
	case "a8":
		address := 0xFF00 + uint16(cpu.ImmediateByte())
		value = uint16(cpu.Gameboy.MMU.ReadByte(address))
	default:
		log.Fatal("unsupported operand named ", operand.Name)
	}

	if !operand.Immediate {
		value = cpu.Gameboy.MMU.ReadWord(value)
	}
	if operand.Increment {
		value += 1
	}
	if operand.Decrement {
		value -= 1
	}

	return value
}

func (cpu *CPU) SetOperand(operand OperandInfo, value uint16) {
	if operand.Increment {
		value += 1
	}
	if operand.Decrement {
		value -= 1
	}
	if !operand.Immediate {
		address := cpu.GetOperand(OperandInfo{Name: operand.Name, Immediate: true})
		cpu.Gameboy.MMU.WriteWord(address, value)
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
	cpu.ProgramCounter += opcode.Bytes
}
