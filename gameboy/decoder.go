package gameboy

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type OperandInfo struct {
	Name      string
	Immediate bool
	Increment bool
	Decrement bool
	Bytes     uint8
	Location  uint16
}

type FlagInfo struct {
	Z string
	N string
	H string
	C string
}

type OpcodeInfo struct {
	Mnemonic  string
	Bytes     uint16
	Cycles    []int
	Operands  []OperandInfo
	Immediate bool
	Flags     FlagInfo
	Hex       string
}

type OpcodesHexMap struct {
	Unprefixed map[string]OpcodeInfo
	Cbprefixed map[string]OpcodeInfo
}

type OpcodesByteMap struct {
	Unprefixed map[byte]OpcodeInfo
	Cbprefixed map[byte]OpcodeInfo
}

type Decoder struct {
	OpcodesMap OpcodesByteMap
}

func ToByteMap(opcodes map[string]OpcodeInfo) map[byte]OpcodeInfo {
	byteMap := make(map[byte]OpcodeInfo)
	for hexOp, info := range opcodes {
		byteOp, _ := strconv.ParseUint(hexOp[2:], 16, 8)
		byteMap[byte(byteOp)] = info
	}

	return byteMap
}

func NewDecoder() Decoder {
	bytes, err := os.ReadFile(filepath.Join(os.Getenv("GOPATH"), "./opcodes.json"))
	if err != nil {
		log.Fatal("Error opening opcodes file: ", err)
	}

	var opcodesHexMap OpcodesHexMap
	json.Unmarshal(bytes, &opcodesHexMap)

	opcodesMap := OpcodesByteMap{
		Unprefixed: ToByteMap(opcodesHexMap.Unprefixed),
		Cbprefixed: ToByteMap(opcodesHexMap.Cbprefixed),
	}

	return Decoder{OpcodesMap: opcodesMap}
}

func (decoder Decoder) DecodeUnprefixed(op byte) OpcodeInfo {
	return decoder.OpcodesMap.Unprefixed[op]
}

func (decoder Decoder) DecodePrefixed(op byte) OpcodeInfo {
	return decoder.OpcodesMap.Cbprefixed[op]
}
