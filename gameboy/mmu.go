package gameboy

import "encoding/binary"

type MMU struct {
	Memory []byte
}

func NewMMU() MMU {
	return MMU{Memory: make([]byte, 2048)}
}

func (mmu MMU) ReadByte(address int) byte {
	return mmu.Memory[address]
}

func (mmu *MMU) WriteByte(address int, value byte) {
	(*mmu).Memory[address] = value
}

func (mmu MMU) ReadWord(address int) uint16 {
	return binary.LittleEndian.Uint16(mmu.Memory[address : address+1])
}

func (mmu *MMU) WriteWord(address int, value uint16) {
	highByte := byte((value & 0xFF00) >> 8)
	lowByte := byte(value & 0x00FF)
	(*mmu).Memory[address] = lowByte
	(*mmu).Memory[address+1] = highByte
}
