package gameboy

import (
	"encoding/binary"
)

type Registers struct {
	A byte
	F Flags
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte
}

type Pair string

const (
	BC Pair = "BC"
	DE Pair = "DE"
	HL Pair = "HL"
	AF Pair = "AF"
)

func (registers Registers) GetPair(pair Pair) uint16 {
	switch pair {
	case AF:
		return binary.LittleEndian.Uint16([]byte{registers.A, registers.F.Byte()})
	case BC:
		return binary.LittleEndian.Uint16([]byte{registers.B, registers.C})
	case DE:
		return binary.LittleEndian.Uint16([]byte{registers.D, registers.E})
	case HL:
		return binary.LittleEndian.Uint16([]byte{registers.H, registers.L})
	}
	return 0
}

func (registers *Registers) SetPair(pair Pair, value uint16) {
	hi := uint8(value >> 8)
	lo := uint8(value & 0xff)

	switch pair {
	case AF:
		registers.A = lo
		registers.F.Set(hi)
	case BC:
		registers.B = lo
		registers.C = hi
	case DE:
		registers.D = lo
		registers.E = hi
	case HL:
		registers.H = lo
		registers.L = hi
	}
}
