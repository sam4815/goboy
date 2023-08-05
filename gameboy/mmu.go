package gameboy

import (
	"encoding/binary"
)

type MMU struct {
	Memory []byte
}

func NewMMU(bytes []byte) MMU {
	memory := make([]byte, 65536)
	copy(memory, bytes)

	memory[0xFF05] = 0x00
	memory[0xFF06] = 0x00
	memory[0xFF07] = 0x00
	memory[0xFF10] = 0x80
	memory[0xFF11] = 0xBF
	memory[0xFF12] = 0xF3
	memory[0xFF14] = 0xBF
	memory[0xFF16] = 0x3F
	memory[0xFF17] = 0x00
	memory[0xFF19] = 0xBF
	memory[0xFF1A] = 0x7F
	memory[0xFF1B] = 0xFF
	memory[0xFF1C] = 0x9F
	memory[0xFF1E] = 0xBF
	memory[0xFF20] = 0xFF
	memory[0xFF21] = 0x00
	memory[0xFF22] = 0x00
	memory[0xFF23] = 0xBF
	memory[0xFF24] = 0x77
	memory[0xFF25] = 0xF3
	memory[0xFF26] = 0xF1
	memory[0xFF40] = 0x91
	memory[0xFF41] = 0x81
	memory[0xFF42] = 0x00
	memory[0xFF43] = 0x00
	memory[0xFF44] = 0x00
	memory[0xFF45] = 0x00
	memory[0xFF47] = 0xFC
	memory[0xFF48] = 0xFF
	memory[0xFF49] = 0xFF
	memory[0xFF4A] = 0x00
	memory[0xFF4B] = 0x00
	memory[0xFFFF] = 0x00

	return MMU{Memory: memory}
}

func (mmu MMU) ReadByte(address uint16) byte {
	return mmu.Memory[address]
}

func (mmu *MMU) WriteByte(address uint16, value byte) {
	(*mmu).Memory[address] = value
}

func (mmu MMU) ReadWord(address uint16) uint16 {
	return binary.LittleEndian.Uint16(mmu.Memory[address : address+2])
}

func (mmu *MMU) WriteWord(address uint16, value uint16) {
	highByte := byte((value & 0xFF00) >> 8)
	lowByte := byte(value & 0x00FF)
	(*mmu).Memory[address] = lowByte
	(*mmu).Memory[address+1] = highByte
}
