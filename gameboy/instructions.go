package gameboy

import (
	"log"
)

func (cpu *CPU) Execute(opcode OpcodeInfo) int {
	// log.Print(opcode.Mnemonic, opcode.Hex)

	switch opcode.Mnemonic {
	case "NOP":
	case "DI":

	case "LD":
		fallthrough
	case "LDH":
		value := cpu.GetOperand(opcode.Operands[1])
		cpu.SetOperand(opcode.Operands[0], value)

	case "INC":
		result := cpu.Registers.F.Add(cpu.GetOperand(opcode.Operands[0]), 1, false, opcode.Flags)
		cpu.SetOperand(opcode.Operands[0], result)

	case "DEC":
		result := cpu.Registers.F.Sub(cpu.GetOperand(opcode.Operands[0]), 1, false, opcode.Flags)
		cpu.SetOperand(opcode.Operands[0], result)

	case "JR":
		if len(opcode.Operands) == 1 || cpu.Registers.F.GetFlagOperand(opcode.Operands[0]) {
			cpu.ProgramCounter += uint16(cpu.ImmediateByteSigned())
		} else {
			return opcode.Cycles[1]
		}

	case "JP":
		if len(opcode.Operands) == 1 || cpu.Registers.F.GetFlagOperand(opcode.Operands[0]) {
			cpu.NextProgramCounter = cpu.ImmediateWord()
		} else {
			return opcode.Cycles[1]
		}

	case "CALL":
		if len(opcode.Operands) == 1 || cpu.Registers.F.GetFlagOperand(opcode.Operands[0]) {
			cpu.PushStack(cpu.StackPointer + 3)
			cpu.NextProgramCounter = cpu.ImmediateWord()
		} else {
			return opcode.Cycles[1]
		}

	case "RET":
		if len(opcode.Operands) == 0 || cpu.Registers.F.GetFlagOperand(opcode.Operands[0]) {
			cpu.NextProgramCounter = cpu.PopStack()
		} else {
			return opcode.Cycles[1]
		}

	case "ADD":
		result := cpu.Registers.F.Add(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), false, opcode.Flags)
		cpu.SetOperand(opcode.Operands[0], result)

	case "ADC":
		result := cpu.Registers.F.Add(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), true, opcode.Flags)
		cpu.SetOperand(opcode.Operands[0], result)

	case "SUB":
		result := cpu.Registers.F.Sub(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), false, opcode.Flags)
		cpu.SetOperand(opcode.Operands[0], result)

	case "SBC":
		result := cpu.Registers.F.Sub(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), true, opcode.Flags)
		cpu.SetOperand(opcode.Operands[0], result)

	case "AND":
		result := cpu.Registers.F.And(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]))
		cpu.SetOperand(opcode.Operands[0], result)

	case "OR":
		result := cpu.Registers.F.Or(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]))
		cpu.SetOperand(opcode.Operands[0], result)

	case "XOR":
		result := cpu.Registers.F.Xor(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]))
		cpu.SetOperand(opcode.Operands[0], result)

	case "CP":
		cpu.Registers.F.Sub(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), false, opcode.Flags)

	case "POP":
		cpu.SetOperand(opcode.Operands[0], cpu.PopStack())

	case "PUSH":
		cpu.PushStack(cpu.GetOperand(opcode.Operands[0]))

	case "RST":
		cpu.PushStack(cpu.ProgramCounter)
		cpu.NextProgramCounter = opcode.Operands[0].Location

	case "SWAP":
		value := cpu.GetOperand(opcode.Operands[0])
		swapped := uint8(value&0xf<<4) | uint8(value>>4)
		cpu.SetOperand(opcode.Operands[0], uint16(swapped))

	case "RRCA":
		value := cpu.Registers.A
		cpu.Registers.F.Carry = value&1 == 1
		cpu.Registers.A = value<<7 | value>>1

	case "RRA":
		cpu.Registers.A = cpu.Registers.F.RotateRightWriteCarry(cpu.Registers.A)

	case "SLA":
		value := cpu.Registers.F.RotateLeftWriteCarry(byte(cpu.GetOperand(opcode.Operands[0])))
		value &= 0
		cpu.SetOperand(opcode.Operands[0], uint16(value))

	case "RL":
		value := cpu.Registers.F.RotateLeftWriteCarry(byte(cpu.GetOperand(opcode.Operands[0])))
		cpu.SetOperand(opcode.Operands[0], uint16(value))

	case "RLC":
		cpu.Registers.F.Carry = cpu.GetOperand(opcode.Operands[0])&0x80 == 0x80
		value := cpu.Registers.F.RotateLeftWriteCarry(byte(cpu.GetOperand(opcode.Operands[0])))
		cpu.SetOperand(opcode.Operands[0], uint16(value))

	default:
		log.Printf("unimplemented opcode %s: %s", opcode.Hex, opcode.Mnemonic)
	}

	return opcode.Cycles[0]
}
