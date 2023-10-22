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
		return binary.BigEndian.Uint16([]byte{registers.A, registers.F.Byte()})
	case BC:
		return binary.BigEndian.Uint16([]byte{registers.B, registers.C})
	case DE:
		return binary.BigEndian.Uint16([]byte{registers.D, registers.E})
	case HL:
		return binary.BigEndian.Uint16([]byte{registers.H, registers.L})
	}
	return 0
}

func (registers *Registers) SetPair(pair Pair, value uint16) {
	lowByte := byte((value & 0xFF00) >> 8)
	highByte := byte(value & 0x00FF)

	switch pair {
	case AF:
		registers.A = lowByte
		registers.F.Set(highByte)
	case BC:
		registers.B = lowByte
		registers.C = highByte
	case DE:
		registers.D = lowByte
		registers.E = highByte
	case HL:
		registers.H = lowByte
		registers.L = highByte
	}
}
