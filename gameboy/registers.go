package gameboy

import (
	"encoding/binary"
	"fmt"
	"log"
)

type Registers struct {
	A byte
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte
	F Flags
}

type Pair uint8

const (
	BC Pair = 0x00
	DE Pair = 0x01
	HL Pair = 0x02
	SP Pair = 0x03
	AF Pair = 0x04
)

func (cpu CPU) GetPair(pair Pair) uint16 {
	switch pair {
	case AF:
		return binary.LittleEndian.Uint16([]byte{cpu.Registers.A, cpu.Registers.F.Byte()})
	case BC:
		return binary.LittleEndian.Uint16([]byte{cpu.Registers.B, cpu.Registers.C})
	case DE:
		return binary.LittleEndian.Uint16([]byte{cpu.Registers.D, cpu.Registers.E})
	case HL:
		return binary.LittleEndian.Uint16([]byte{cpu.Registers.H, cpu.Registers.L})
	case SP:
		return cpu.StackPointer
	}

	err := fmt.Errorf("attempted to access invalid pair %d", pair)
	log.Fatal(err)

	return 0
}

func (cpu *CPU) SetPair(pair Pair, value uint16) {
	hi := uint8(value >> 8)
	lo := uint8(value & 0xff)

	switch pair {
	case AF:
		cpu.Registers.A = lo
		cpu.Registers.F.Set(hi)
	case BC:
		cpu.Registers.B = lo
		cpu.Registers.C = hi
	case DE:
		cpu.Registers.D = lo
		cpu.Registers.E = hi
	case HL:
		cpu.Registers.H = lo
		cpu.Registers.L = hi
	case SP:
		cpu.StackPointer = value
	}
}
