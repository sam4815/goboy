package gameboy

import (
	"log"
)

func (cpu *CPU) Execute(opcode OpcodeInfo) int {
	// fmt.Printf("Executing %s %s\n", opcode.Mnemonic, opcode.Hex)
	switch opcode.Mnemonic {
	case "NOP":
	case "STOP":
	case "HALT":

	case "LD":
		fallthrough
	case "LDH":
		value := cpu.GetOperand(opcode.Operands[1])
		cpu.SetOperand(opcode.Operands[0], value)

	case "LDHLSP":
		immediate := cpu.ImmediateByte()
		result, flags := Add(cpu.GetOperand(opcode.Operands[1]), uint16(immediate), false)
		if immediate > 127 {
			result, flags = Sub(cpu.GetOperand(opcode.Operands[0]), uint16(immediate-127), false)
		}

		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "INC":
		result, flags := Add(cpu.GetOperand(opcode.Operands[0]), 1, false)
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "DEC":
		result, flags := Sub(cpu.GetOperand(opcode.Operands[0]), 1, false)
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "JR":
		if len(opcode.Operands) == 1 {
			cpu.NextProgramCounter += cpu.GetOperand(opcode.Operands[0])
		} else if cpu.Registers.F.GetFlagOperand(opcode.Operands[0]) {
			cpu.ProgramCounter += cpu.GetOperand(opcode.Operands[1])
		} else {
			return opcode.Cycles[1]
		}

	case "JP":
		if len(opcode.Operands) == 1 {
			cpu.NextProgramCounter = cpu.GetOperand(opcode.Operands[0])
		} else if cpu.Registers.F.GetFlagOperand(opcode.Operands[0]) {
			cpu.NextProgramCounter = cpu.GetOperand(opcode.Operands[1])
		} else {
			return opcode.Cycles[1]
		}

	case "CALL":
		if len(opcode.Operands) == 1 || cpu.Registers.F.GetFlagOperand(opcode.Operands[0]) {
			cpu.PushStack(cpu.ProgramCounter + 3)
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
		result, flags := Add(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), false)
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "ADD16":
		result, flags := Add16(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), false)
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "ADC":
		result, flags := Add(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), cpu.Registers.F.Carry)
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "SUB":
		result, flags := Sub(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), false)
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "SBC":
		result, flags := Sub(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), cpu.Registers.F.Carry)
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "AND":
		result, flags := And(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]))
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "DAA":
		value := cpu.Registers.A
		carry := cpu.Registers.F.Carry

		if !cpu.Registers.F.Subtract {
			if cpu.Registers.F.Carry || value > 0x99 {
				carry = true
				cpu.Registers.A += 0x60
			}
			if cpu.Registers.F.HalfCarry || (value&0x0f) > 0x09 {
				cpu.Registers.A += 0x6
			}
		} else {
			if cpu.Registers.F.Carry {
				cpu.Registers.A -= 0x60
			}
			if cpu.Registers.F.HalfCarry {
				cpu.Registers.A -= 0x6
			}
		}

		cpu.Registers.F.ProcessFlags(Flags{Zero: cpu.Registers.A == 0, Carry: carry}, opcode.Flags)

	case "OR":
		result, flags := Or(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]))
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "XOR":
		result, flags := Xor(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]))
		cpu.SetOperand(opcode.Operands[0], result)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "CPL":
		cpu.Registers.A ^= 0xff
		cpu.Registers.F.ProcessFlags(Flags{}, opcode.Flags)

	case "CP":
		_, flags := Sub(cpu.GetOperand(opcode.Operands[0]), cpu.GetOperand(opcode.Operands[1]), false)
		cpu.Registers.F.ProcessFlags(flags, opcode.Flags)

	case "POP":
		cpu.SetOperand(opcode.Operands[0], cpu.PopStack())

	case "PUSH":
		cpu.PushStack(cpu.GetOperand(opcode.Operands[0]))

	case "RST":
		cpu.PushStack(cpu.ProgramCounter + 1)
		cpu.NextProgramCounter = opcode.Operands[0].Location

	case "RETI":
		cpu.NextProgramCounter = cpu.PopStack()

	case "SWAP":
		value := cpu.GetOperand(opcode.Operands[0])
		swapped := uint8(value&0xf<<4) | uint8(value>>4)
		cpu.SetOperand(opcode.Operands[0], uint16(swapped))

	case "RRCA":
		value := cpu.Registers.A
		cpu.Registers.A = value<<7 | value>>1
		cpu.Registers.F.ProcessFlags(Flags{Carry: value&1 == 1}, opcode.Flags)

	case "RRA":
		value := cpu.Registers.A
		cpu.Registers.A = BoolToByte(cpu.Registers.F.Carry)<<7 | value>>1
		cpu.Registers.F.ProcessFlags(Flags{Carry: value&1 == 1}, opcode.Flags)

	case "SLA":
		value := byte(cpu.GetOperand(opcode.Operands[0]))
		cpu.SetOperand(opcode.Operands[0], uint16(value<<1))
		cpu.Registers.F.ProcessFlags(Flags{Zero: value<<1 == 0, Carry: value>>7 == 1}, opcode.Flags)

	case "RL":
		value := byte(cpu.GetOperand(opcode.Operands[0]))
		result := (value << 1) | BoolToByte(cpu.Registers.F.Carry)
		cpu.SetOperand(opcode.Operands[0], uint16(result))
		cpu.Registers.F.ProcessFlags(Flags{Zero: value<<1 == 0, Carry: value>>7 == 1}, opcode.Flags)

	case "RLA":
		value := cpu.Registers.A
		cpu.Registers.A = (value << 1) | BoolToByte(cpu.Registers.F.Carry)
		cpu.Registers.F.ProcessFlags(Flags{Carry: value>>7 == 1}, opcode.Flags)

	case "RLC":
		value := cpu.GetOperand(opcode.Operands[0])
		result := value<<1 | value>>7
		cpu.SetOperand(opcode.Operands[0], uint16(result))
		cpu.Registers.F.ProcessFlags(Flags{Zero: result == 0, Carry: value>>7 == 1}, opcode.Flags)

	case "RLCA":
		value := cpu.Registers.A
		cpu.Registers.A = value<<1 | value>>7
		cpu.Registers.F.ProcessFlags(Flags{Carry: value>>7 == 1}, opcode.Flags)

	case "SCF":
		cpu.Registers.F.ProcessFlags(Flags{}, opcode.Flags)

	case "CCF":
		cpu.Registers.F.ProcessFlags(Flags{Carry: !cpu.Registers.F.Carry}, opcode.Flags)

	case "DI":
		cpu.EnableInterrupts = false

	case "EI":
		cpu.EnableInterrupts = true

	default:
		log.Printf("unimplemented opcode %s: %s", opcode.Hex, opcode.Mnemonic)
	}

	return opcode.Cycles[0]
}
