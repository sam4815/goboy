package gameboy

import (
	"encoding/binary"
	"log"
)

type CPU struct {
	ProgramCounter byte
	StackPointer   byte
	Registers      Registers
}

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

type Flags struct {
	Zero      bool
	Subtract  bool
	HalfCarry bool
	Carry     bool
}

func (flags Flags) Byte() byte {
	flagByte := 0

	if flags.Zero {
		flagByte |= (1 << 7)
	}
	if flags.Subtract {
		flagByte |= (1 << 6)
	}
	if flags.HalfCarry {
		flagByte |= (1 << 5)
	}
	if flags.Carry {
		flagByte |= (1 << 4)
	}

	return byte(flagByte)
}

func (registers Registers) AF() uint16 {
	return binary.LittleEndian.Uint16([]byte{registers.A, registers.F.Byte()})
}

func (registers Registers) BC() uint16 {
	return binary.LittleEndian.Uint16([]byte{registers.B, registers.C})
}

func (registers Registers) DE() uint16 {
	return binary.LittleEndian.Uint16([]byte{registers.D, registers.E})
}

func (registers Registers) HL() uint16 {
	return binary.LittleEndian.Uint16([]byte{registers.H, registers.L})
}

func (cpu *CPU) Step(memory *[]byte) {
	instruction := (*memory)[cpu.ProgramCounter]

	log.Print(instruction)
}
