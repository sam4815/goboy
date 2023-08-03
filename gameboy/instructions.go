package gameboy

import (
	"fmt"
	"log"
)

func (cpu *CPU) Execute(opcode OpcodeInfo) int {
	if opcode.Mnemonic == "NOP" {
		return opcode.Cycles[0]
	}

	if opcode.Mnemonic == "DI" {
		return opcode.Cycles[0]
	}

	if opcode.Mnemonic == "LD" || opcode.Mnemonic == "LDH" {
		value := cpu.GetOperand(opcode.Operands[1])
		cpu.SetOperand(opcode.Operands[0], value)

		return opcode.Cycles[0]
	}

	err := fmt.Errorf("unimplemented opcode %s: %s", opcode.Hex, opcode.Mnemonic)
	log.Print(err)

	return 0
}
