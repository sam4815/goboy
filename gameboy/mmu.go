package gameboy

import "encoding/binary"

type MMU struct {
	Memory []byte
}

func NewMMU(bytes []byte) MMU {
	memory := make([]byte, 65536)
	copy(memory, bytes)

	return MMU{Memory: memory}
}

func (mmu MMU) ReadByte(address uint16) byte {
	return mmu.Memory[address]
}

func (mmu *MMU) WriteByte(address uint16, value byte) {
	(*mmu).Memory[address] = value
}

func (mmu MMU) ReadWord(address uint16) uint16 {
	return binary.LittleEndian.Uint16(mmu.Memory[address : address+1])
}

func (mmu *MMU) WriteWord(address uint16, value uint16) {
	highByte := byte((value & 0xFF00) >> 8)
	lowByte := byte(value & 0x00FF)
	(*mmu).Memory[address] = lowByte
	(*mmu).Memory[address+1] = highByte
}
