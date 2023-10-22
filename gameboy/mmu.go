package gameboy

import (
	"log"
)

type BankingMode uint8

const (
	None BankingMode = 0
	MBC1 BankingMode = 1
	MBC2 BankingMode = 2
)

type MBC1Registers struct {
	ROMBank byte
	RAMBank byte
	RAMOn   bool
	Mode    byte
}

type MMU struct {
	Gameboy       *Gameboy
	ROM           []byte
	ExternalRAM   []byte
	WorkRAM       []byte
	HighRAM       []byte
	OAM           []byte
	ROMOffset     uint16
	RAMOffset     uint16
	Keys          byte
	Interrupt     byte
	BankingMode   BankingMode
	MBC1Registers MBC1Registers
}

func (mmu *MMU) InitializeIO() {
	// mmu.Memory[0xFF05] = 0x00
	// mmu.Memory[0xFF06] = 0x00
	// mmu.Memory[0xFF07] = 0x00
	// mmu.Memory[0xFF10] = 0x80
	// mmu.Memory[0xFF11] = 0xBF
	// mmu.Memory[0xFF12] = 0xF3
	// mmu.Memory[0xFF14] = 0xBF
	// mmu.Memory[0xFF16] = 0x3F
	// mmu.Memory[0xFF17] = 0x00
	// mmu.Memory[0xFF19] = 0xBF
	// mmu.Memory[0xFF1A] = 0x7F
	// mmu.Memory[0xFF1B] = 0xFF
	// mmu.Memory[0xFF1C] = 0x9F
	// mmu.Memory[0xFF1E] = 0xBF
	// mmu.Memory[0xFF20] = 0xFF
	// mmu.Memory[0xFF21] = 0x00
	// mmu.Memory[0xFF22] = 0x00
	// mmu.Memory[0xFF23] = 0xBF
	// mmu.Memory[0xFF24] = 0x77
	// mmu.Memory[0xFF25] = 0xF3
	// mmu.Memory[0xFF26] = 0xF1
	// mmu.Memory[0xFF40] = 0x91
	// mmu.Memory[0xFF41] = 0x81
	// mmu.Memory[0xFF42] = 0x00
	// mmu.Memory[0xFF43] = 0x00
	// mmu.Memory[0xFF44] = 0x00
	// mmu.Memory[0xFF45] = 0x00
	// mmu.Memory[0xFF47] = 0xFC
	// mmu.Memory[0xFF48] = 0xFF
	// mmu.Memory[0xFF49] = 0xFF
	// mmu.Memory[0xFF4A] = 0x00
	// mmu.Memory[0xFF4B] = 0x00
	// mmu.Memory[0xFFFF] = 0x00
}

func (mmu *MMU) ReadBankingMode() {
	switch mmu.ROM[0x147] {
	case 1:
		mmu.BankingMode = MBC1
	case 2:
		mmu.BankingMode = MBC1
	case 3:
		mmu.BankingMode = MBC1
	case 5:
		mmu.BankingMode = MBC2
	case 6:
		mmu.BankingMode = MBC2
	}
}

func NewMMU(bytes []byte) *MMU {
	mmu := MMU{
		ROM:         make([]byte, 16384),
		ExternalRAM: make([]byte, 8192),
		WorkRAM:     make([]byte, 8192),
		HighRAM:     make([]byte, 8192),
		OAM:         make([]byte, 160),
		ROMOffset:   0x0000,
	}
	copy(mmu.ROM, bytes)

	mmu.InitializeIO()
	mmu.ReadBankingMode()

	return &mmu
}

func (mmu MMU) ReadByte(address uint16) byte {
	switch address & 0xF000 {
	case 0x0000, 0x1000, 0x2000, 0x3000:
		return mmu.ROM[address]

	case 0x4000, 0x5000, 0x6000, 0x7000:
		return mmu.ROM[mmu.ROMOffset+(address&0x3FFF)]

	case 0x8000, 0x9000:
		return mmu.Gameboy.GPU.VRAM[address&0x1FFF]

	case 0xA000, 0xB000:
		return mmu.ExternalRAM[mmu.RAMOffset+address&0x1FFF]

	case 0xC000, 0xD000, 0xE000:
		return mmu.WorkRAM[address&0x1FFF]

	case 0xF000:
		switch address & 0xF00 {
		case 0x100, 0x200, 0x300, 0x400, 0x500, 0x600, 0x700, 0x800, 0x900, 0xA00, 0xB00, 0xC00, 0xD00:
			return mmu.WorkRAM[address&0x1FFF]

		case 0xE00:
			return mmu.OAM[address&0xFF]

		case 0xF00:
			switch address & 0xF0 {
			case 0x00, 0x10, 0x20, 0x30:
				switch address & 0xF {
				case 0x00:
					return mmu.Keys
				}
			case 0x40, 0x50, 0x60, 0x70:
				return mmu.Gameboy.GPU.ReadRegisters(address)
			case 0x80, 0x90, 0xA0, 0xB0, 0xC0, 0xD0, 0xE0:
				return mmu.HighRAM[address&0x7F]
			case 0xF0:
				switch address & 0xF {
				case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE:
					return mmu.HighRAM[address&0x7F]
				case 0xF:
					return mmu.Interrupt
				}
			}
		}
	}

	log.Printf("attempting to access invalid memory address %x", address)
	return 0
}

func (mmu *MMU) WriteByte(address uint16, value byte) {
	// if value != 0 {
	// 	log.Print("MMU WRITE BYTE: ", address, value)
	// 	mmu.Gameboy.Dump()
	// }
	switch address & 0xF000 {
	case 0x0000, 0x1000:
		switch mmu.BankingMode {
		case 2, 3:
			mmu.MBC1Registers.RAMOn = (value & 0x0F) == 0x0A
		}

	case 0x2000, 0x3000:
		switch mmu.BankingMode {
		case 1, 2, 3:
			lower5Bits := value & 0x1F
			if lower5Bits == 0 {
				lower5Bits = 1
			}

			mmu.MBC1Registers.ROMBank &= 0x60
			mmu.MBC1Registers.ROMBank |= lower5Bits

			mmu.ROMOffset = uint16(mmu.MBC1Registers.ROMBank) * 0x4000
		}

	case 0x4000, 0x5000:
		switch mmu.BankingMode {
		case 1, 2, 3:
			if mmu.MBC1Registers.Mode == 1 {
				mmu.MBC1Registers.RAMBank = value & 0x03
				mmu.RAMOffset = uint16(mmu.MBC1Registers.RAMBank) * 0x2000
			} else {
				upper2Bits := value & 0x03 << 5

				mmu.MBC1Registers.ROMBank &= 0x1F
				mmu.MBC1Registers.ROMBank |= upper2Bits
				mmu.ROMOffset = uint16(mmu.MBC1Registers.ROMBank) * 0x4000
			}
		}

	case 0x6000, 0x7000:
		switch mmu.BankingMode {
		case 2, 3:
			mmu.MBC1Registers.Mode = value & 0x01
		}

	case 0x8000, 0x9000:
		mmu.Gameboy.GPU.VRAM[address&0x1FFF] = value
		mmu.Gameboy.GPU.UpdateTile(address&0x1FFF, value)

	case 0xA000, 0xB000:
		mmu.ExternalRAM[mmu.RAMOffset+address&0x1FFF] = value

	case 0xC000, 0xD000, 0xE000:
		mmu.WorkRAM[address&0x1FFF] = value

	case 0xF000:
		switch address & 0xF00 {
		case 0x100, 0x200, 0x300, 0x400, 0x500, 0x600, 0x700, 0x800, 0x900, 0xA00, 0xB00, 0xC00, 0xD00:
			mmu.WorkRAM[address&0x1FFF] = value

		case 0xF00:
			switch address & 0xF0 {
			case 0x00, 0x10, 0x20, 0x30:
				switch address & 0xF {
				case 0x00:
					mmu.Keys = value
				}
			case 0x40, 0x50, 0x60, 0x70:
				mmu.Gameboy.GPU.WriteRegisters(address, value)
			case 0x80, 0x90, 0xA0, 0xB0, 0xC0, 0xD0, 0xE0:
				mmu.HighRAM[address&0x7F] = value
			case 0xF0:
				switch address & 0xF {
				case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE:
					mmu.HighRAM[address&0x7F] = value
				case 0xF:
					mmu.Interrupt = value
				}
			}
		}
	}
}

func (mmu MMU) ReadWord(address uint16) uint16 {
	return uint16(mmu.ReadByte(address)) + (uint16(mmu.ReadByte(address+1)) << 8)
}

func (mmu *MMU) WriteWord(address uint16, value uint16) {
	// if value != 0 {
	// 	log.Print("WRITING WORD ", address, value, byte((value&0xFF00)>>8), byte(value&0x00FF))
	// }
	highByte := byte((value & 0xFF00) >> 8)
	lowByte := byte(value & 0x00FF)

	mmu.WriteByte(address, lowByte)
	mmu.WriteByte(address+1, highByte)
}
