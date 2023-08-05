package gameboy

type Key string

const (
	Start       Key = "Start"
	Select      Key = "Select"
	A           Key = "A"
	B           Key = "B"
	Up          Key = "Up"
	Down        Key = "Down"
	Left        Key = "Left"
	Right       Key = "Right"
	Unsupported Key = "Unsupported"
)

func (mmu *MMU) KeyDown(key Key) {
	switch key {
	case Right, A:
		mmu.Keys &= 0b1110
	case Left, B:
		mmu.Keys &= 0b1101
	case Up, Select:
		mmu.Keys &= 0b1011
	case Down, Start:
		mmu.Keys &= 0b0111
	}
}

func (mmu *MMU) KeyUp(key Key) {
	switch key {
	case Right, A:
		mmu.Keys |= 0b1110
	case Left, B:
		mmu.Keys |= 0b1101
	case Up, Select:
		mmu.Keys |= 0b1011
	case Down, Start:
		mmu.Keys |= 0b0111
	}
}

func (mmu *MMU) WriteKeys(b byte) {
	mmu.Keys &= 0b00110000
}
