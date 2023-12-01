package gameboy

func Add(a uint16, b uint16, carry bool) (uint16, Flags) {
	result := int16(a) + int16(b)

	if carry {
		result += 1
		b -= 1
	}

	return uint16(result), Flags{
		Zero:      result == 0,
		Subtract:  false,
		Carry:     result > 0xFF,
		HalfCarry: (a&0x0F)+(b&0x0F) > 0x0F,
	}
}

func Add16(a uint16, b uint16, carry bool) (uint16, Flags) {
	result := int32(a) + int32(b)

	if carry {
		result += 1
		b -= 1
	}

	return uint16(result), Flags{
		Zero:      result == 0,
		Subtract:  false,
		Carry:     result > 0xFFFF,
		HalfCarry: (a&0x0FFF)+(b&0x0FFF) > 0x0FFF,
	}
}

func Sub16(a uint16, b uint16, carry bool) (uint16, Flags) {
	result := int32(a) - int32(b)

	if carry {
		result -= 1
		a -= 1
	}

	return uint16(result), Flags{
		Zero:      result == 0,
		Subtract:  true,
		Carry:     result < 0,
		HalfCarry: (a & 0x0FFF) < (b & 0x0FFF),
	}
}

func Sub(a uint16, b uint16, carry bool) (uint16, Flags) {
	result := int16(a) - int16(b)

	if carry {
		result -= 1
		a -= 1
	}

	return uint16(result), Flags{
		Zero:      result == 0,
		Subtract:  true,
		Carry:     result < 0,
		HalfCarry: (a & 0x0F) < (b & 0x0F),
	}
}

func And(a uint16, b uint16) (uint16, Flags) {
	result := a & b

	return result, Flags{Zero: result == 0}
}

func Or(a uint16, b uint16) (uint16, Flags) {
	result := a | b

	return result, Flags{Zero: result == 0}
}

func Xor(a uint16, b uint16) (uint16, Flags) {
	result := a ^ b

	return result, Flags{Zero: result == 0}
}

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
