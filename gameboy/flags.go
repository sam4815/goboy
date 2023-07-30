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
