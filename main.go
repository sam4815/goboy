package main

import (
	"flag"
	"goboy/gameboy"
	"log"
	"os"
)

func main() {
	filePath := flag.String("file", "", "gameboy rom file path")
	flag.Parse()

	bytes, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}

	gb := gameboy.Gameboy{Memory: bytes, CPU: gameboy.CPU{}}

	gb.Run()
}
