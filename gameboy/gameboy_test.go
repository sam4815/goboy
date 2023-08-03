package gameboy

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("GOPATH", "../")
}

func TestTitle(t *testing.T) {
	bytes, err := os.ReadFile("../snake.gb")
	if err != nil {
		t.Errorf("Error opening file: %s", err)
	}

	gameboy := New(bytes)

	if gameboy.Metadata.Title != "Yvar's GB Snake" {
		t.Errorf("expected Yvar's GB Snake, got %s", gameboy.Metadata.Title)
	}
}

func TestCPU10(t *testing.T) {
	bytes, err := os.ReadFile("../snake.gb")
	if err != nil {
		t.Errorf("Error opening file: %s", err)
	}

	gameboy := New(bytes)

	for i := 0; i < 5; i++ {
		gameboy.CPU.Step()
	}

	expectedProgramCounter := uint16(345)
	if gameboy.CPU.ProgramCounter != expectedProgramCounter {
		t.Errorf("expected %d, got %d", expectedProgramCounter, gameboy.CPU.ProgramCounter)
	}

	expectedStackPointer := uint16(65534)
	if gameboy.CPU.StackPointer != expectedStackPointer {
		t.Errorf("expected %d, got %d", expectedStackPointer, gameboy.CPU.StackPointer)
	}

	expectedRegisters := Registers{
		A: 17,
		F: Flags{
			Zero:      true,
			Subtract:  false,
			Carry:     false,
			HalfCarry: false,
		},
		B: 128,
		C: 0,
		D: 255,
		E: 86,
		H: 0,
		L: 13,
	}

	if gameboy.CPU.Registers != expectedRegisters {
		t.Errorf("expected %v, got %v", expectedRegisters, gameboy.CPU.Registers)
	}
}
