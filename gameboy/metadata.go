package gameboy

import "encoding/binary"

type Metadata struct {
	Title           string
	CGB             byte
	NewLicenseeCode uint16
	SGB             byte
	CartridgeType   byte
	RomSize         byte
	RamSize         byte
	DestinationCode byte
	OldLicenseeCode byte
	MaskRomVersion  byte
	HeaderChecksum  byte
	GlobalChecksum  uint16
}

func ParseMetadata(bytes []byte) Metadata {
	return Metadata{
		Title:           string(bytes[0x134:0x143]),
		CGB:             bytes[0x143],
		NewLicenseeCode: binary.LittleEndian.Uint16(bytes[0x144:0x146]),
		SGB:             bytes[0x146],
		CartridgeType:   bytes[0x147],
		RomSize:         bytes[0x148],
		RamSize:         bytes[0x149],
		DestinationCode: bytes[0x14A],
		OldLicenseeCode: bytes[0x14B],
		MaskRomVersion:  bytes[0x14C],
		HeaderChecksum:  bytes[0x14D],
		GlobalChecksum:  binary.LittleEndian.Uint16(bytes[0x14E:0x150]),
	}
}
