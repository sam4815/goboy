package gameboy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type CPUState struct {
	Pc  uint16
	Sp  uint16
	A   byte
	B   byte
	C   byte
	D   byte
	E   byte
	F   byte
	H   byte
	L   byte
	Ram [][]int
}

type TestCase struct {
	Name    string
	Initial CPUState
	Final   CPUState
}

func init() {
	os.Setenv("GOPATH", "../")
}

func (state CPUState) Print() string {
	ramStrs := make([]string, 0)
	for _, ramPiece := range state.Ram {
		ramStrs = append(ramStrs, fmt.Sprintf("(%d, %d)", ramPiece[0], ramPiece[1]))
	}

	return fmt.Sprint(
		" PC: ", state.Pc,
		", SP: ", state.Sp,
		", A: ", state.A,
		", B: ", state.B,
		", C: ", state.C,
		", D: ", state.D,
		", E: ", state.E,
		", F: ", state.F,
		", H: ", state.H,
		", L: ", state.L,
		", RAM: ", strings.Join(ramStrs, ", "),
	)
}

func (state *CPUState) RemoveInvalidRAM() CPUState {
	withoutShadow := make([][]int, 0)
	for _, ramPiece := range state.Ram {
		if ramPiece[0] >= 0xE000 && ramPiece[0] <= 0xFDFF && ramPiece[1] != 0 {
			withoutShadow = append(withoutShadow, []int{ramPiece[0] - 0x2000, ramPiece[1]})
		} else if ramPiece[1] != 0 {
			withoutShadow = append(withoutShadow, ramPiece)
		}
	}

	state.Ram = withoutShadow

	return *state
}

func copyState(gameboy *Gameboy, state CPUState) {
	gameboy.CPU.ProgramCounter = state.Pc
	gameboy.CPU.StackPointer = state.Sp
	gameboy.CPU.Registers.A = state.A
	gameboy.CPU.Registers.B = state.B
	gameboy.CPU.Registers.C = state.C
	gameboy.CPU.Registers.D = state.D
	gameboy.CPU.Registers.E = state.E
	gameboy.CPU.Registers.F.Set(state.F)
	gameboy.CPU.Registers.H = state.H
	gameboy.CPU.Registers.L = state.L

	for _, memoryPiece := range state.Ram {
		if memoryPiece[0] < 16384*2 {
			gameboy.MMU.ROM[memoryPiece[0]] = byte(memoryPiece[1])
		} else {
			gameboy.MMU.WriteByte(uint16(memoryPiece[0]), byte(memoryPiece[1]))
		}
	}
}

func dumpState(gameboy Gameboy) CPUState {
	ram := make([][]int, 0)
	for i, memory := range gameboy.MMU.ROM {
		if memory != 0 {
			ram = append(ram, []int{i, int(memory)})
		}
	}
	for i, memory := range gameboy.GPU.VRAM {
		if memory != 0 {
			ram = append(ram, []int{i + 0x8000, int(memory)})
		}
	}
	for i, memory := range gameboy.MMU.ExternalRAM {
		if memory != 0 {
			ram = append(ram, []int{i + 0xA000, int(memory)})
		}
	}
	for i, memory := range gameboy.MMU.WorkRAM {
		if memory != 0 {
			ram = append(ram, []int{i + 0xC000, int(memory)})
			// ram = append(ram, []int{i + 0xE000, int(memory)})
		}
	}
	for i, memory := range gameboy.MMU.HighRAM {
		if memory != 0 {
			ram = append(ram, []int{i + 0xFF80, int(memory)})
		}
	}

	return CPUState{
		Pc:  gameboy.CPU.ProgramCounter,
		Sp:  gameboy.CPU.StackPointer,
		A:   gameboy.CPU.Registers.A,
		B:   gameboy.CPU.Registers.B,
		C:   gameboy.CPU.Registers.C,
		D:   gameboy.CPU.Registers.D,
		E:   gameboy.CPU.Registers.E,
		F:   gameboy.CPU.Registers.F.Byte(),
		H:   gameboy.CPU.Registers.H,
		L:   gameboy.CPU.Registers.L,
		Ram: ram,
	}
}

func TestCPU(t *testing.T) {
	for i := range make([]int, 256) {
		filename := fmt.Sprintf("./gameboy/v1/%02x.json", i)
		bytes, err := os.ReadFile(filepath.Join(os.Getenv("GOPATH"), filename))
		if err != nil {
			log.Fatal("Error opening opcodes file: ", err)
		}

		var testCases []TestCase
		json.Unmarshal(bytes, &testCases)

		failed := 0

		for _, testCase := range testCases {
			gameboy := Gameboy{
				MMU: NewMMU(make([]byte, 16384*128)),
				CPU: NewCPU(),
				GPU: NewGPU(),
			}
			gameboy.CPU.Gameboy = &gameboy
			gameboy.GPU.Gameboy = &gameboy
			gameboy.MMU.Gameboy = &gameboy
			gameboy.MMU.BankingMode = MBC1
			gameboy.MMU.WriteByte(0x2000, 1)

			copyState(&gameboy, testCase.Initial)

			gameboy.CPU.Step()

			dumped := dumpState(gameboy)
			final := testCase.Final.RemoveInvalidRAM()

			if !reflect.DeepEqual(dumped, final) {
				// fmt.Printf("ERROR for test case %s\n\033[32mexpected %s\033[0m \n\033[31mreceived %s\033[0m\n", testCase.Name, final.Print(), dumped.Print())
				failed += 1
			}
		}

		fmt.Printf("0x%02x \033[32m%d passed\033[0m \033[31m%d failed\033[0m\n", i, len(testCases)-failed, failed)
	}
}
