package gameboy

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

func (flags *Flags) Set(b byte) {
	flags.Zero = b>>7 == 1
	flags.Subtract = b>>6 == 1
	flags.HalfCarry = b>>5 == 1
	flags.Carry = b>>4 == 1
}

func (flags *Flags) Add(a uint16, b uint16, carry bool) uint16 {
	result := int16(a) + int16(b)

	if carry && flags.Carry {
		result += 1
		b -= 1
	}

	flags.Zero = result == 0
	flags.Subtract = false
	flags.Carry = result > 0xFF
	flags.HalfCarry = (a&0x0F)+(b&0x0F) > 0x0F

	return uint16(result)
}

func (flags *Flags) Sub(a uint16, b uint16, carry bool) uint16 {
	result := int16(a) - int16(b)

	if carry && flags.Carry {
		result -= 1
		b -= 1
	}

	flags.Zero = result == 0
	flags.Subtract = true
	flags.Carry = result < 0
	flags.HalfCarry = (a & 0x0F) < (b & 0x0F)

	return uint16(result)
}

func (flags *Flags) And(a uint16, b uint16) uint16 {
	result := a & b

	flags.Zero = result == 0
	flags.Subtract = false
	flags.Carry = false
	flags.HalfCarry = true

	return result
}

func (flags *Flags) Or(a uint16, b uint16) uint16 {
	result := a | b

	flags.Zero = result == 0
	flags.Subtract = false
	flags.Carry = false
	flags.HalfCarry = false

	return result
}

func (flags *Flags) Xor(a uint16, b uint16) uint16 {
	result := a ^ b

	flags.Zero = result == 0
	flags.Subtract = false
	flags.Carry = false
	flags.HalfCarry = false

	return result
}

func (flags Flags) GetFlagOperand(operand OperandInfo) bool {
	switch operand.Name {
	case "Z":
		return flags.Zero
	case "NZ":
		return !flags.Zero
	case "C":
		return flags.Carry
	case "NC":
		return !flags.Carry
	}
	return false
}

func (flags *Flags) RotateLeftWriteCarry(n byte) byte {
	rotated := flags.RotateLeftReadCarry(n)
	flags.Carry = n>>7 == 0x01
	return rotated
}

func (flags Flags) RotateLeftReadCarry(n byte) byte {
	rotated := n << 1
	if flags.Carry {
		rotated |= 0x01
	}
	return rotated
}

func (flags *Flags) RotateRightWriteCarry(n byte) byte {
	rotated := flags.RotateRightReadCarry(n)
	flags.Carry = n&0x01 == 0x01
	return rotated
}

func (flags Flags) RotateRightReadCarry(n byte) byte {
	rotated := n >> 1
	if flags.Carry {
		rotated |= 0x80
	}
	return rotated
}
