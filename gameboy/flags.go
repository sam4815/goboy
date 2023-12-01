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
	flags.Zero = (b>>7)&1 == 1
	flags.Subtract = (b>>6)&1 == 1
	flags.HalfCarry = (b>>5)&1 == 1
	flags.Carry = (b>>4)&1 == 1
}

func (flags *Flags) ProcessFlags(result Flags, instructions FlagInstructions) {
	switch instructions.Z {
	case "Z":
		flags.Zero = result.Zero
	case "1":
		flags.Zero = true
	case "0":
		flags.Zero = false
	}

	switch instructions.N {
	case "N":
		flags.Subtract = result.Subtract
	case "1":
		flags.Subtract = true
	case "0":
		flags.Subtract = false
	}

	switch instructions.C {
	case "C":
		flags.Carry = result.Carry
	case "1":
		flags.Carry = true
	case "0":
		flags.Carry = false
	}

	switch instructions.H {
	case "H":
		flags.HalfCarry = result.HalfCarry
	case "1":
		flags.HalfCarry = true
	case "0":
		flags.HalfCarry = false
	}
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

// func (flags *Flags) RotateLeftWriteCarry(n byte) byte {
// 	rotated := flags.RotateLeftReadCarry(n)
// 	flags.Carry = n>>7 == 0x01
// 	return rotated
// }

// func (flags Flags) RotateLeftReadCarry(n byte) byte {
// 	rotated := n << 1
// 	if flags.Carry {
// 		rotated |= 0x01
// 	}
// 	return rotated
// }

// func (flags *Flags) RotateRightWriteCarry(n byte) byte {
// 	rotated := flags.RotateRightReadCarry(n)
// 	flags.Carry = n&0x01 == 0x01
// 	return rotated
// }

// func (flags Flags) RotateRightReadCarry(n byte) byte {
// 	rotated := n >> 1
// 	if flags.Carry {
// 		rotated |= 0x80
// 	}
// 	return rotated
// }
