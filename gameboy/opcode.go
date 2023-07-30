package gameboy

import (
	"fmt"
	"log"
)

type Opcode struct {
	Hex   uint8
	HiHex uint8
	LoHex uint8
}

func DecodeOpcode(b byte) Opcode {
	return Opcode{
		HiHex: b / 16,
		LoHex: b % 16,
		Hex:   b,
	}
}

func (cpu *CPU) Execute(opcode Opcode) (int, int) {
	if opcode.Hex == 0x0 {
		return 1, 4
	}

	if opcode.LoHex == 0x1 && opcode.HiHex <= 0x3 {
		cpu.SetPair(Pair(opcode.HiHex), cpu.Immediate16())

		return 3, 12
	}

	if opcode.LoHex == 0x2 && opcode.HiHex <= 0x1 {
		address := cpu.GetPair(Pair(opcode.HiHex))
		cpu.Gameboy.MMU.WriteByte(address, cpu.Registers.A)

		return 1, 8
	}

	if opcode.LoHex == 0x3 && opcode.HiHex <= 0x3 {
		cpu.SetPair(Pair(opcode.HiHex), cpu.GetPair(Pair(opcode.HiHex))+1)

		return 1, 8
	}

	err := fmt.Errorf("unimplemented opcode 0x%x", opcode.Hex)
	log.Fatal(err)

	return 0, 0
}
