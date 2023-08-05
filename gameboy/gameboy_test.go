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

func TestCPU100(t *testing.T) {
	bytes, err := os.ReadFile("../snake.gb")
	if err != nil {
		t.Errorf("Error opening file: %s", err)
	}

	gameboy := New(bytes)

	for i := 0; i < 100; i++ {
		gameboy.CPU.Step()
	}

	expectedProgramCounter := uint16(349)
	if gameboy.CPU.ProgramCounter != expectedProgramCounter {
		t.Errorf("expected %d, got %d", expectedProgramCounter, gameboy.CPU.ProgramCounter)
	}

	expectedStackPointer := uint16(65534)
	if gameboy.CPU.StackPointer != expectedStackPointer {
		t.Errorf("expected %d, got %d", expectedStackPointer, gameboy.CPU.StackPointer)
	}

	expectedCycles := 1068
	if gameboy.CPU.Cycles != expectedCycles {
		t.Errorf("expected %d, got %d", expectedCycles, gameboy.CPU.Cycles)
	}

	expectedRegisters := Registers{
		A: 0,
		F: Flags{
			Zero:      false,
			Subtract:  true,
			Carry:     true,
			HalfCarry: false,
		},
		B: 128,
		C: 0,
		D: 0,
		E: 216,
		H: 1,
		L: 77,
	}

	if gameboy.CPU.Registers != expectedRegisters {
		t.Errorf("expected %v, got %v", expectedRegisters, gameboy.CPU.Registers)
	}
}
