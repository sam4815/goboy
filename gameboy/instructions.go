package gameboy

import (
	"fmt"
	"log"
)

func (cpu *CPU) Execute(opcode OpcodeInfo) int {
	if opcode.Mnemonic == "NOP" {
		return 0
	}

	// if opcode.LoHex == 0x1 && opcode.HiHex <= 0x3 {
	// 	cpu.SetPair(Pair(opcode.HiHex), cpu.ImmediateWord())

	// 	return 3, 12
	// }

	// if opcode.LoHex == 0x2 && opcode.HiHex <= 0x1 {
	// 	address := cpu.GetPair(Pair(opcode.HiHex))
	// 	cpu.Gameboy.MMU.WriteByte(address, cpu.Registers.A)

	// 	return 1, 8
	// }

	// if opcode.LoHex == 0x3 && opcode.HiHex <= 0x3 {
	// 	cpu.SetPair(Pair(opcode.HiHex), cpu.GetPair(Pair(opcode.HiHex))+1)

	// 	return 1, 8
	// }

	err := fmt.Errorf("unimplemented opcode %s: %s", opcode.Hex, opcode.Mnemonic)
	log.Print(err)

	return 0
}
